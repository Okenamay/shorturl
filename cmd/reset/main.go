package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// 1. Константы и глобальные переменные
const (
	resetMarker = "// generate:reset"
	genFileName = "reset.gen.go"
)

// funcTpl - это шаблон для генерации метода Reset(), он использует флаги из
// fieldInfo для построения корректной логики
var funcTpl = `
// Reset сбрасывает состояние объекта к начальным значениям.
func ({{.Receiver}} *{{.TypeName}}) Reset() {
	if {{.Receiver}} == nil {
		return
	}

{{range .Fields}}
	{{- if .IsPrimitive}}
	// Сброс примитива (int, string, bool)
	{{.Receiver}}.{{.Name}} = *new({{.TypeName}})
	{{- else if .IsSlice}}
	// Сброс слайса
	{{.Receiver}}.{{.Name}} = {{.Receiver}}.{{.Name}}[:0]
	{{- else if .IsMap}}
	// Очистка мапы
	clear({{.Receiver}}.{{.Name}})
	{{- else if .IsPointerToPrimitive}}
	// Сброс значения по указателю на примитив
	if {{.Receiver}}.{{.Name}} != nil {
		*{{.Receiver}}.{{.Name}} = *new({{.TypeName}})
	}
	{{- else if .IsPointerToResettable}}
	// Вызов Reset() для указателя на структуру
	if {{.Receiver}}.{{.Name}} != nil {
		{{.Receiver}}.{{.Name}}.Reset()
	}
	{{- else if .IsResettableValue}}
	// Вызов Reset() для вложенной структуры (по значению)
	(&{{.Receiver}}.{{.Name}}).Reset()
	{{- else if .IsPointerToNonResettable}}
	// Указатель на структуру без Reset() - ничего не делаем
	{{- end}}
{{end}}
}
`

// tmpl - скомпилированный шаблон
var tmpl = template.Must(template.New("reset").Parse(funcTpl))

// 2. Структуры данных для шаблона

// fieldInfo хранит всю информацию о поле, извлеченную из go/types. Эта
// информация будет передана в шаблон
type fieldInfo struct {
	Name     string // Имя поля - "i", "strP", "child" и т.п.
	TypeName string // Имя типа - "int", "string", "ResetableStruct" и т.п.

	// Флаги, определяющие логику сброса
	IsPrimitive              bool // int, string, bool...
	IsSlice                  bool // []T
	IsMap                    bool // map[K]V
	IsPointerToPrimitive     bool // *string, *int
	IsPointerToResettable    bool // *MyStruct (где MyStruct имеет Reset())
	IsResettableValue        bool // MyStruct (где *MyStruct имеет Reset())
	IsPointerToNonResettable bool // *MyStruct (где MyStruct НЕ имеет Reset())
}

// structInfo хранит всю информацию о найденной структуре
type structInfo struct {
	PkgName  string      // Имя пакета, например, "config"
	TypeName string      // Имя структуры, например, "Cfg"
	Receiver string      // Имя получателя для метода, например, "c"
	Fields   []fieldInfo // Список проанализированных полей
}

// 3. Точка входа
func main() {
	log.Println("Starting reset generator (go/types + template)...")

	// 1) Загружаем пакеты из текущей директории (./...)
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:   ".",   // Сканируем из корня проекта
		Tests: false, // Игнорируем тестовые файлы
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Fatalf("Failed to load packages: %v", err)
	}

	// map[pkg_dir_path] -> []structInfo
	pkgStructs := make(map[string][]structInfo)
	pkgNames := make(map[string]string)

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			for _, e := range pkg.Errors {
				log.Printf("Skipping package %s due to error: %v", pkg.PkgPath, e)
			}
			continue
		}

		var pkgDir string

		// 2) Итерируем по файлам пакета, ищем структуры
		for i, fileAST := range pkg.Syntax {
			if len(pkg.GoFiles) > i {
				pkgDir = filepath.Dir(pkg.GoFiles[i])
				if strings.HasSuffix(pkg.GoFiles[i], genFileName) {
					continue
				}
			}

			// 3) Ищем структуры с маркером
			ast.Inspect(fileAST, func(n ast.Node) bool {
				genDecl, ok := n.(*ast.GenDecl)
				if !ok {
					return true
				}
				if !hasResetMarker(genDecl.Doc) {
					return true
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					// 4) Нашли, получаем информацию о типе из go/types
					obj := pkg.TypesInfo.Defs[typeSpec.Name]
					if obj == nil {
						log.Printf("Warning: could not find type info for %s", typeSpec.Name.Name)
						continue
					}

					structType, ok := obj.Type().Underlying().(*types.Struct)
					if !ok {
						continue // Маркер стоит не на структуре
					}

					log.Printf("Found marker on struct %s in %s", typeSpec.Name.Name, pkg.PkgPath)

					// 5) Собираем информацию о структуре для шаблона
					si := structInfo{
						PkgName:  pkg.Name,
						TypeName: typeSpec.Name.Name,
						Receiver: strings.ToLower(string(typeSpec.Name.Name[0])),
						Fields:   parseStructFields(structType),
					}

					pkgStructs[pkgDir] = append(pkgStructs[pkgDir], si)
					pkgNames[pkgDir] = pkg.Name
				}
				return true
			})
		}
	}

	// 6) Генерируем и записываем файлы
	for pkgDir, structs := range pkgStructs {
		if len(structs) > 0 {
			err := generateFile(pkgDir, pkgNames[pkgDir], structs)
			if err != nil {
				log.Printf("Failed to generate file for %s: %v", pkgDir, err)
			} else {
				log.Printf("Successfully generated %s in %s", genFileName, pkgDir)
			}
		}
	}

	log.Println("Reset generator finished.")
}

// 4. Анализ пакетов и типов:
// hasResetMarker проверяет, содержит ли группа комментариев маркер
func hasResetMarker(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	for _, comment := range doc.List {
		if strings.TrimSpace(comment.Text) == resetMarker {
			return true
		}
	}
	return false
}

// parseStructFields анализирует *types.Struct и возвращает срез fieldInfo
func parseStructFields(s *types.Struct) []fieldInfo {
	var fields []fieldInfo
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		if field.Embedded() {
			continue // Пропускаем вложенные анонимные поля
		}

		fi := fieldInfo{
			Name: field.Name(),
		}

		// Заполняем флаги в fi, анализируя тип поля
		analyzeFieldType(&fi, field.Type())
		fields = append(fields, fi)
	}
	return fields
}

// analyzeFieldType рекурсивно анализирует types.Type и заполняет fieldInfo
func analyzeFieldType(fi *fieldInfo, typ types.Type) {
	// 1) Получаем имя типа (например, "string", "MyStruct" и т.п.)
	// types.TypeString(typ) может вернуть "main.MyStruct", нам это не нужно,
	// мы ищем именованный тип.
	if named, ok := typ.(*types.Named); ok {
		fi.TypeName = named.Obj().Name()
	}

	// 2) Рекурсивно "разворачиваем" указатели
	if ptr, ok := typ.(*types.Pointer); ok {
		// Рекурсивно вызываем для типа, на который он указывает (например, string)
		analyzeFieldType(fi, ptr.Elem())

		// После рекурсии, fi.TypeName будет "string", а fi.IsPrimitive = true,
		// необходимо скорректировать флаги для указателя.
		if fi.IsPrimitive {
			fi.IsPrimitive = false
			fi.IsPointerToPrimitive = true
		} else if fi.IsResettableValue {
			fi.IsResettableValue = false
			fi.IsPointerToResettable = true
		} else {
			fi.IsPointerToNonResettable = true
		}
		return
	}

	// 3) Анализируем базовый (underlying) тип
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		fi.IsPrimitive = true
		if fi.TypeName == "" { // Если это не type MyInt int, а просто int
			fi.TypeName = t.Name()
		}
	case *types.Slice:
		fi.IsSlice = true
	case *types.Map:
		fi.IsMap = true
	case *types.Named, *types.Struct:
		// 4) Вложенные структуры (не указатели)
		// Проверяем, есть ли у *T (указателя на структуру) метод Reset()
		ptrType := types.NewPointer(typ)
		if hasResetMethod(ptrType) {
			fi.IsResettableValue = true
		}
		if fi.TypeName == "" && t.String() == "struct{}" {
			// Игнорируем анонимные пустые структуры
			return
		}
	default:
		// Другие типы (chan, func, interface) игнорируем
	}
}

// hasResetMethod проверяет, реализует ли данный тип interface{ Reset() }
// (с поинтер-ресивером)
func hasResetMethod(typ types.Type) bool {
	// Ищем метод Reset()
	obj, _, _ := types.LookupFieldOrMethod(typ, true, nil, "Reset")

	if obj == nil {
		return false
	}

	meth, ok := obj.(*types.Func)
	if !ok {
		return false
	}

	// Проверяем сигнатуру: func()
	sig, ok := meth.Type().(*types.Signature)
	if !ok {
		return false
	}

	// Убеждаемся, что это поинтер-ресивер (или интерфейс) и что у него нет
	// аргументов и нет возвращаемых значений
	recv := sig.Recv()

	if recv == nil {
		return false
	}

	_, isPtr := recv.Type().(*types.Pointer)

	isPtrRecv := types.IsInterface(recv.Type()) || isPtr

	return isPtrRecv && sig.Params().Len() == 0 && sig.Results().Len() == 0
}

// 5. Генерация файла

// generateFile Генерирует и записывает файл, используя шаблон
func generateFile(dir string, pkgName string, structs []structInfo) error {
	var buf bytes.Buffer

	// 1) Заголовок файла
	fmt.Fprintf(&buf, "// Code generated by cmd/reset. DO NOT EDIT.\n\npackage %s\n\n", pkgName)

	// 2) Генерируем методы для каждой структуры
	for _, s := range structs {
		// tmpl.Execute записывает результат в буфер
		err := tmpl.Execute(&buf, s)
		if err != nil {
			return fmt.Errorf("ошибка выполнения шаблона для %s: %w", s.TypeName, err)
		}
	}

	// 3) Форматируем сгенерированный код (go fmt)
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// Если `go fmt` упал, это значит, мы сгенерировали
		// невалидный код. Логируем ошибку И сам код.
		log.Printf("--- ОШИБКА ФОРМАТИРОВАНИЯ ---\n%s\n----------------------------", buf.String())
		return fmt.Errorf("не удалось отформатировать код для %s: %w", dir, err)
	}

	// 4) Записываем в файл
	outputPath := filepath.Join(dir, genFileName)
	err = os.WriteFile(outputPath, formatted, 0644)
	if err != nil {
		return fmt.Errorf("не удалось записать %s: %w", outputPath, err)
	}

	return nil
}
