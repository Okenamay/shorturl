#!/bin/bash

# --- НАСТРОЙКА ---
# Укажите директорию, где лежат ваши .pprof файлы.
# "." означает текущую директорию.
PROFILE_DIR="."
# -----------------

# Эта опция (nullglob) не запускает цикл, если файлы не найдены,
# вместо того чтобы пытаться обработать '*.result.pprof' как имя файла.
shopt -s nullglob

echo "Starting pprof comparison in: $PROFILE_DIR"

# 1. Ищем все .result.pprof файлы в указанной директории
for result_file in "$PROFILE_DIR"/*.result.pprof; do
    
    # 2. Формируем имя парного .base.pprof файла
    # (Заменяем суффикс ".result.pprof" на ".base.pprof")
    base_file="${result_file/.result.pprof/.base.pprof}"
    
    # 3. Проверяем, существует ли парный .base.pprof файл
    if [ -f "$base_file" ]; then
        echo "=================================================="
        echo "Processing pair:"
        echo "  Base:   $base_file"
        echo "  Result: $result_file"
        echo "--------------------------------------------------"
        
        # 4. Запускаем саму команду pprof
        # (Вы можете добавить сюда другие флаги, например, -http=:8080)
        # ИСПРАВЛЕНО: Используем 'go tool pprof', чтобы не зависеть от $PATH
        go tool pprof -top -diff_base="$base_file" "$result_file"
        
        echo # Добавляем пустую строку для читаемости
    else
        echo "=================================================="
        echo "Skipping: $result_file"
        echo "  Reason: Base file not found at: $base_file"
    fi
done

echo "=================================================="
echo "Comparison finished."

# Выключаем опцию, чтобы не влиять на другие скрипты
shopt -u nullglob

