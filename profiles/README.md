$ ./compare_pprof.sh
Starting pprof comparison in: .
==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_app_checker.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_app_checker.result.pprof
--------------------------------------------------
File: checker.test.exe
Build ID: C:\Users\Yamabaka\Dev\shorturl\checker.test.exe2025-11-01 22:24:04.2380457 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:30:45 MSK
Showing nodes accounting for 0.24GB, 5.70% of 4.27GB total
Dropped 16 nodes (cum <= 0.02GB)
      flat  flat%   sum%        cum   cum%
    0.17GB  4.01%  4.01%     0.19GB  4.50%  net/url.parse
    0.06GB  1.46%  5.47%     0.25GB  5.96%  net/url.ParseRequestURI
    0.02GB  0.49%  5.96%     0.02GB  0.49%  errors.New (inline)
   -0.01GB  0.26%  5.70%     0.24GB  5.70%  github.com/Okenamay/shorturl.git/internal/app/checker.CheckURL
         0     0%  5.70%     0.19GB  4.48%  github.com/Okenamay/shorturl.git/internal/app/checker.BenchmarkCheckURL_InvalidParse
         0     0%  5.70%     0.11GB  2.67%  github.com/Okenamay/shorturl.git/internal/app/checker.BenchmarkCheckURL_InvalidScheme
         0     0%  5.70%    -0.06GB  1.44%  github.com/Okenamay/shorturl.git/internal/app/checker.BenchmarkCheckURL_Valid
         0     0%  5.70%     0.24GB  5.70%  testing.(*B).launch
         0     0%  5.70%     0.24GB  5.70%  testing.(*B).runN

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_app_hasher.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_app_hasher.result.pprof
--------------------------------------------------
File: hasher.test.exe
Build ID: C:\Users\Yamabaka\Dev\shorturl\hasher.test.exe2025-11-01 22:24:06.3495447 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:30:48 MSK
Showing nodes accounting for -815.06MB, 54.88% of 1485.10MB total
Dropped 27 nodes (cum <= 7.43MB)
      flat  flat%   sum%        cum   cum%
 -582.55MB 39.23% 39.23%  -582.55MB 39.23%  crypto/md5.New (inline)
 -388.51MB 26.16% 65.39%  -388.51MB 26.16%  encoding/hex.EncodeToString (inline)
     226MB 15.22% 50.17%  -814.56MB 54.85%  github.com/Okenamay/shorturl.git/internal/app/hasher.ShortenURL
  -94.50MB  6.36% 56.53%   -94.50MB  6.36%  crypto/md5.(*digest).Sum (inline)
   24.50MB  1.65% 54.88%    24.50MB  1.65%  io.WriteString
         0     0% 54.88%  -814.56MB 54.85%  github.com/Okenamay/shorturl.git/internal/app/hasher.BenchmarkShortenURL
         0     0% 54.88%  -814.56MB 54.85%  testing.(*B).launch
         0     0% 54.88%  -814.56MB 54.85%  testing.(*B).runN

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_app_middleware_auth.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_app_middleware_auth.result.pprof
--------------------------------------------------
File: auth.test.exe
Build ID: C:\Users\Yamabaka\Dev\shorturl\auth.test.exe2025-11-01 22:10:45.8930681 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:30:55 MSK
Showing nodes accounting for 41.02MB, 0.93% of 4410.63MB total
Dropped 1 node (cum <= 22.05MB)
      flat  flat%   sum%        cum   cum%
     -25MB  0.57%  0.57%      -25MB  0.57%  strings.(*Builder).grow
     -19MB  0.43%     1%        3MB 0.068%  github.com/golang-jwt/jwt/v4.(*SigningMethodHMAC).Verify
  -17.51MB   0.4%  1.39%      -14MB  0.32%  net/http.readRequest
  -17.50MB   0.4%  1.79%     2.50MB 0.057%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.buildJWTString
      15MB  0.34%  1.45%       15MB  0.34%  reflect.copyVal
  -13.50MB  0.31%  1.76%   -13.50MB  0.31%  github.com/golang-jwt/jwt/v4.splitToken
   13.50MB  0.31%  1.45%    13.50MB  0.31%  crypto/internal/fips140/sha256.(*Digest).Sum
   12.50MB  0.28%  1.17%    46.50MB  1.05%  github.com/golang-jwt/jwt/v4.(*Parser).ParseUnverified
   12.50MB  0.28%  0.88%    12.50MB  0.28%  encoding/base64.(*Encoding).DecodeString
      11MB  0.25%  0.64%    22.50MB  0.51%  encoding/json.Unmarshal
      11MB  0.25%  0.39%    15.50MB  0.35%  crypto/internal/fips140/hmac.New[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }]
      10MB  0.23%  0.16%       10MB  0.23%  net/http.(*Request).WithContext (partial-inline)
      10MB  0.23% 0.068%       38MB  0.86%  github.com/golang-jwt/jwt/v4.(*Parser).ParseWithClaims
    9.02MB   0.2%  0.27%     9.02MB   0.2%  bufio.NewReaderSize (inline)
    7.50MB  0.17%  0.44%       17MB  0.39%  github.com/golang-jwt/jwt/v4.(*NumericDate).UnmarshalJSON
      -7MB  0.16%  0.28%       -7MB  0.16%  encoding/json.(*Decoder).refill
       7MB  0.16%  0.44%        7MB  0.16%  encoding/json.NewDecoder (inline)
      -7MB  0.16%  0.28%       -7MB  0.16%  net/textproto.MIMEHeader.Add (inline)
       7MB  0.16%  0.44%        7MB  0.16%  net/url.parse
      -6MB  0.14%  0.31%       -6MB  0.14%  net/http.Header.Clone (inline)
      -6MB  0.14%  0.17%       -6MB  0.14%  github.com/golang-jwt/jwt/v4.RegisteredClaims.Valid
       6MB  0.14%  0.31%       17MB  0.39%  encoding/json.(*decodeState).literalStore
   -5.50MB  0.12%  0.18%    -5.50MB  0.12%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.getUserID.func1
       5MB  0.11%  0.29%        5MB  0.11%  net/http/httptest.NewRecorder (inline)
      -5MB  0.11%  0.18%       -5MB  0.11%  github.com/golang-jwt/jwt/v4.NewParser (inline)
    4.50MB   0.1%  0.28%     4.50MB   0.1%  crypto/internal/fips140/sha256.New (inline)
    4.50MB   0.1%  0.39%    18.51MB  0.42%  encoding/json.Marshal
    4.50MB   0.1%  0.49%     4.50MB   0.1%  github.com/golang-jwt/jwt/v4.NewNumericDate (inline)
    4.50MB   0.1%  0.59%     4.50MB   0.1%  github.com/google/uuid.NewRandomFromReader
       4MB 0.091%  0.68%    29.50MB  0.67%  encoding/json.(*decodeState).object
       4MB 0.091%  0.77%        4MB 0.091%  github.com/golang-jwt/jwt/v4.NewWithClaims (inline)
   -3.50MB 0.079%  0.69%    -2.50MB 0.057%  github.com/golang-jwt/jwt/v4.NumericDate.MarshalJSON
    2.51MB 0.057%  0.75%     2.51MB 0.057%  sync.(*Pool).pinSlow
    2.50MB 0.057%  0.81%     2.50MB 0.057%  net/http.readCookies
    2.50MB 0.057%  0.86%     2.50MB 0.057%  reflect.mapassign_faststr0
   -2.50MB 0.057%  0.81%    -2.50MB 0.057%  encoding/base64.(*Encoding).EncodeToString
       2MB 0.045%  0.85%        2MB 0.045%  strings.NewReader (inline)
      -2MB 0.045%  0.81%       -2MB 0.045%  net/textproto.(*Reader).ReadLine (inline)
    1.50MB 0.034%  0.84%        2MB 0.045%  github.com/golang-jwt/jwt/v4.(*SigningMethodHMAC).Sign
   -1.50MB 0.034%  0.81%    -1.50MB 0.034%  net/textproto.readMIMEHeader
    1.50MB 0.034%  0.84%    -3.48MB 0.079%  net/http/httptest.NewRequestWithContext
    1.50MB 0.034%  0.87%     1.50MB 0.034%  reflect.New
       1MB 0.023%   0.9%        1MB 0.023%  context.WithValue
       1MB 0.023%  0.92%        1MB 0.023%  strconv.FormatFloat (inline)
       1MB 0.023%  0.94%        1MB 0.023%  encoding/json.(*scanner).pushParseState
    0.50MB 0.011%  0.95%     0.50MB 0.011%  runtime.allocm
   -0.50MB 0.011%  0.94%    -0.50MB 0.011%  internal/sync.newIndirectNode[go.shape.interface {},go.shape.interface {}] (inline)
    0.50MB 0.011%  0.95%    33.50MB  0.76%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.getUserID
   -0.50MB 0.011%  0.94%    -0.50MB 0.011%  runtime.acquireSudog
   -0.50MB 0.011%  0.93%       15MB  0.34%  encoding/json.mapEncoder.encode
    0.50MB 0.011%  0.94%     0.50MB 0.011%  bytes.(*Buffer).grow
    0.50MB 0.011%  0.95%     0.50MB 0.011%  github.com/google/uuid.UUID.String
   -0.50MB 0.011%  0.94%    -0.50MB 0.011%  net/textproto.NewReader (inline)
   -0.50MB 0.011%  0.93%    -0.50MB 0.011%  sync.OnceFunc (inline)
    0.50MB 0.011%  0.94%     0.50MB 0.011%  crypto/internal/fips140/ecdsa.init
   -0.50MB 0.011%  0.93%    -0.50MB 0.011%  unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1
    0.50MB 0.011%  0.94%    -8.51MB  0.19%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.BenchmarkAuthenticatorMiddleware.BenchmarkAuthenticatorMiddleware.Authenticator.func4.func5
   -0.50MB 0.011%  0.93%    -0.50MB 0.011%  runtime.(*timers).addHeap
         0     0%  0.93%     9.02MB   0.2%  bufio.NewReader (inline)
         0     0%  0.93%     0.50MB 0.011%  bytes.(*Buffer).WriteByte
         0     0%  0.93%     4.50MB   0.1%  crypto.Hash.New
         0     0%  0.93%    15.50MB  0.35%  crypto/hmac.New
         0     0%  0.93%     4.50MB   0.1%  crypto/hmac.New.UnwrapNew[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }].func1
         0     0%  0.93%    -0.50MB 0.011%  crypto/internal/fips140/ed25519.init
         0     0%  0.93%    13.50MB  0.31%  crypto/internal/fips140/hmac.(*HMAC).Sum
         0     0%  0.93%     4.50MB   0.1%  crypto/sha256.New
         0     0%  0.93%    13.50MB  0.31%  encoding/json.(*Decoder).Decode
         0     0%  0.93%       -4MB 0.091%  encoding/json.(*Decoder).readValue
         0     0%  0.93%       -1MB 0.023%  encoding/json.(*decodeState).scanWhile
         0     0%  0.93%    28.50MB  0.65%  encoding/json.(*decodeState).unmarshal
         0     0%  0.93%    29.50MB  0.67%  encoding/json.(*decodeState).value
         0     0%  0.93%    10.50MB  0.24%  encoding/json.(*encodeState).marshal
         0     0%  0.93%    10.50MB  0.24%  encoding/json.(*encodeState).reflectValue
         0     0%  0.93%    -1.50MB 0.034%  encoding/json.appendCompact
         0     0%  0.93%    -0.50MB 0.011%  encoding/json.cachedTypeFields
         0     0%  0.93%       -1MB 0.023%  encoding/json.checkValid
         0     0%  0.93%    -4.50MB   0.1%  encoding/json.indirect
         0     0%  0.93%       -4MB 0.091%  encoding/json.marshalerEncoder
         0     0%  0.93%     2.01MB 0.045%  encoding/json.newEncodeState
         0     0%  0.93%    -1.50MB 0.034%  encoding/json.newScanner
         0     0%  0.93%    -0.50MB 0.011%  encoding/json.newStructEncoder (inline)
         0     0%  0.93%    -0.50MB 0.011%  encoding/json.newTypeEncoder
         0     0%  0.93%        1MB 0.023%  encoding/json.stateBeginValue
         0     0%  0.93%       -4MB 0.091%  encoding/json.structEncoder.encode
         0     0%  0.93%    -0.50MB 0.011%  encoding/json.typeEncoder
         0     0%  0.93%    -0.50MB 0.011%  encoding/json.valueEncoder
         0     0%  0.93%       -6MB  0.14%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.BenchmarkAuthenticatorMiddleware.func1 
         0     0%  0.93%   -60.49MB  1.37%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.BenchmarkAuthenticatorMiddleware.func2 
         0     0%  0.93%    53.50MB  1.21%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.BenchmarkAuthenticatorMiddleware.func3 
         0     0%  0.93%    42.50MB  0.96%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.BenchmarkBuildJWTString
         0     0%  0.93%     7.01MB  0.16%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.BenchmarkGetUserID
         0     0%  0.93%    -0.50MB 0.011%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.init
         0     0%  0.93%       11MB  0.25%  github.com/golang-jwt/jwt/v4.(*Token).SignedString
         0     0%  0.93%    17.01MB  0.39%  github.com/golang-jwt/jwt/v4.(*Token).SigningString
         0     0%  0.93%    12.50MB  0.28%  github.com/golang-jwt/jwt/v4.DecodeSegment
         0     0%  0.93%    -2.50MB 0.057%  github.com/golang-jwt/jwt/v4.EncodeSegment (inline)
         0     0%  0.93%       33MB  0.75%  github.com/golang-jwt/jwt/v4.ParseWithClaims
         0     0%  0.93%    -0.50MB 0.011%  github.com/golang-jwt/jwt/v4.newNumericDateFromSeconds
         0     0%  0.93%     4.50MB   0.1%  github.com/google/uuid.New
         0     0%  0.93%     4.50MB   0.1%  github.com/google/uuid.NewRandom
         0     0%  0.93%    -0.50MB 0.011%  internal/sync.(*HashTrieMap[go.shape.interface {},go.shape.interface {}]).Load
         0     0%  0.93%    -0.50MB 0.011%  internal/sync.(*HashTrieMap[go.shape.interface {},go.shape.interface {}]).init (inline)
         0     0%  0.93%    -0.50MB 0.011%  internal/sync.(*HashTrieMap[go.shape.interface {},go.shape.interface {}]).initSlow
         0     0%  0.93%    -3.50MB 0.079%  net/http.(*Cookie).String
         0     0%  0.93%     2.50MB 0.057%  net/http.(*Request).Cookie
         0     0%  0.93%    -8.51MB  0.19%  net/http.HandlerFunc.ServeHTTP (partial-inline)
         0     0%  0.93%       -7MB  0.16%  net/http.Header.Add (inline)
         0     0%  0.93%      -14MB  0.32%  net/http.ReadRequest
         0     0%  0.93%   -10.50MB  0.24%  net/http.SetCookie
         0     0%  0.93%       -6MB  0.14%  net/http/httptest.(*ResponseRecorder).WriteHeader
         0     0%  0.93%    -3.48MB 0.079%  net/http/httptest.NewRequest (inline)
         0     0%  0.93%    -1.50MB 0.034%  net/textproto.(*Reader).ReadMIMEHeader (inline)
         0     0%  0.93%        7MB  0.16%  net/url.ParseRequestURI
         0     0%  0.93%        8MB  0.18%  reflect.(*MapIter).Key
         0     0%  0.93%        7MB  0.16%  reflect.(*MapIter).Value
         0     0%  0.93%     2.50MB 0.057%  reflect.Value.SetMapIndex
         0     0%  0.93%     2.50MB 0.057%  reflect.mapassign_faststr
         0     0%  0.93%    -0.50MB 0.011%  runtime.(*scavengerState).sleep
         0     0%  0.93%    -0.50MB 0.011%  runtime.(*timer).maybeAdd
         0     0%  0.93%    -0.50MB 0.011%  runtime.(*timer).modify
         0     0%  0.93%    -0.50MB 0.011%  runtime.(*timer).reset (inline)
         0     0%  0.93%    -0.50MB 0.011%  runtime.bgscavenge
         0     0%  0.93%    -0.50MB 0.011%  runtime.chanrecv
         0     0%  0.93%    -0.50MB 0.011%  runtime.chanrecv1
         0     0%  0.93%    -0.50MB 0.011%  runtime.doInit (inline)
         0     0%  0.93%    -0.50MB 0.011%  runtime.doInit1
         0     0%  0.93%    -0.50MB 0.011%  runtime.main
         0     0%  0.93%     0.50MB 0.011%  runtime.mcall
         0     0%  0.93%     0.50MB 0.011%  runtime.newm
         0     0%  0.93%     0.50MB 0.011%  runtime.park_m
         0     0%  0.93%     0.50MB 0.011%  runtime.resetspinning
         0     0%  0.93%     0.50MB 0.011%  runtime.schedule
         0     0%  0.93%     0.50MB 0.011%  runtime.startm
         0     0%  0.93%       -1MB 0.023%  runtime.unique_runtime_registerUniqueMapCleanup.func2
         0     0%  0.93%     0.50MB 0.011%  runtime.wakep
         0     0%  0.93%      -25MB  0.57%  strings.(*Builder).Grow
         0     0%  0.93%   -21.50MB  0.49%  strings.Join
         0     0%  0.93%    -0.50MB 0.011%  sync.(*Map).Load (inline)
         0     0%  0.93%        1MB 0.023%  sync.(*Pool).Get
         0     0%  0.93%     1.50MB 0.034%  sync.(*Pool).Put
         0     0%  0.93%     2.51MB 0.057%  sync.(*Pool).pin
         0     0%  0.93%    42.52MB  0.96%  testing.(*B).launch
         0     0%  0.93%    42.52MB  0.96%  testing.(*B).runN
         0     0%  0.93%    -0.50MB 0.011%  unique.registerCleanup.func1

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_app_middleware_gzipper.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_app_middleware_gzipper.result.pprof
--------------------------------------------------
File: gzipper.test.exe
Build ID: C:\Users\Yamabaka\Dev\shorturl\gzipper.test.exe2025-11-01 22:10:49.7288265 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:31:00 MSK
Showing nodes accounting for -12759.36MB, 66.50% of 19185.82MB total
Dropped 70 nodes (cum <= 95.93MB)
      flat  flat%   sum%        cum   cum%
-4757.09MB 24.79% 24.79% -8608.95MB 44.87%  compress/flate.NewWriter (inline)
-2922.97MB 15.24% 40.03% -2922.97MB 15.24%  compress/flate.(*dictDecoder).init (inline)
-2374.10MB 12.37% 52.40% -3851.86MB 20.08%  compress/flate.(*compressor).init
-1451.74MB  7.57% 59.97% -1451.74MB  7.57%  compress/flate.newDeflateFast (inline)
 -755.39MB  3.94% 63.91%  -755.39MB  3.94%  bufio.NewReaderSize (inline)
 -674.09MB  3.51% 67.42% -3597.06MB 18.75%  compress/flate.NewReader
 -134.59MB   0.7% 68.12% -9198.05MB 47.94%  compress/gzip.NewReader
  130.04MB  0.68% 67.45%   130.04MB  0.68%  net/http.Header.Clone (inline)
   83.56MB  0.44% 67.01%    83.56MB  0.44%  github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper.init.func2
   75.02MB  0.39% 66.62%    75.02MB  0.39%  net/textproto.MIMEHeader.Set (inline)
      15MB 0.078% 66.54%   228.09MB  1.19%  github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper.BenchmarkCompressorMiddleware.func1 
       4MB 0.021% 66.52% -8332.76MB 43.43%  github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper.BenchmarkCompressorMiddleware.BenchmarkCompressorMiddleware.Compressor.func2.func3
       3MB 0.016% 66.50% -4016.60MB 20.94%  github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper.BenchmarkDecompressorMiddleware.BenchmarkDecompressorMiddleware.Decompressor.func2.func3
         0     0% 66.50%  -755.39MB  3.94%  bufio.NewReader (inline)
         0     0% 66.50% -3947.42MB 20.57%  compress/gzip.(*Reader).Reset
         0     0% 66.50% -3597.06MB 18.75%  compress/gzip.(*Reader).readHeader
         0     0% 66.50% -8456.91MB 44.08%  compress/gzip.(*Writer).Write
         0     0% 66.50%   214.51MB  1.12%  github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper.(*gzipWriter).Write
         0     0% 66.50% -8278.76MB 43.15%  github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper.BenchmarkCompressorMiddleware       
         0     0% 66.50% -4563.19MB 23.78%  github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper.BenchmarkDecompressorMiddleware     
         0     0% 66.50% -12350.67MB 64.37%  net/http.HandlerFunc.ServeHTTP (partial-inline)
         0     0% 66.50%    75.02MB  0.39%  net/http.Header.Set (partial-inline)
         0     0% 66.50%   186.54MB  0.97%  net/http/httptest.(*ResponseRecorder).Write
         0     0% 66.50%   130.04MB  0.68%  net/http/httptest.(*ResponseRecorder).WriteHeader
         0     0% 66.50%   133.54MB   0.7%  net/http/httptest.(*ResponseRecorder).writeHeader
         0     0% 66.50%  -500.08MB  2.61%  net/http/httptest.NewRequest (inline)
         0     0% 66.50%  -500.08MB  2.61%  net/http/httptest.NewRequestWithContext
         0     0% 66.50% -12841.32MB 66.93%  testing.(*B).launch
         0     0% 66.50% -12841.95MB 66.93%  testing.(*B).runN

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_app_urlmaker.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_app_urlmaker.result.pprof
--------------------------------------------------
File: urlmaker.test.exe
Build ID: C:\Users\Yamabaka\Dev\shorturl\urlmaker.test.exe2025-11-01 22:10:53.1074553 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:31:04 MSK
Showing nodes accounting for -675.05MB, 54.90% of 1229.59MB total
Dropped 30 nodes (cum <= 6.15MB)
      flat  flat%   sum%        cum   cum%
 -438.04MB 35.62% 35.62%  -438.04MB 35.62%  crypto/md5.New (inline)
 -303.01MB 24.64% 60.27%  -303.01MB 24.64%  encoding/hex.EncodeToString (inline)
     125MB 10.17% 50.10%  -674.04MB 54.82%  github.com/Okenamay/shorturl.git/internal/app/hasher.ShortenURL
  -67.50MB  5.49% 55.59%   -67.50MB  5.49%  crypto/md5.(*digest).Sum (inline)
    8.50MB  0.69% 54.90%     8.50MB  0.69%  io.WriteString
         0     0% 54.90%  -675.54MB 54.94%  github.com/Okenamay/shorturl.git/internal/app/urlmaker.BenchmarkProcessURL
         0     0% 54.90%  -675.54MB 54.94%  github.com/Okenamay/shorturl.git/internal/app/urlmaker.ProcessURL
         0     0% 54.90%  -675.54MB 54.94%  testing.(*B).launch
         0     0% 54.90%  -675.54MB 54.94%  testing.(*B).runN

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_audit.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_audit.result.pprof
--------------------------------------------------
File: audit.test.exe
Build ID: C:\Users\Yamabaka\Dev\shorturl\audit.test.exe2025-11-01 22:11:03.0979649 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:31:11 MSK
Showing nodes accounting for 1109.13MB, 88.46% of 1253.77MB total
Dropped 103 nodes (cum <= 6.27MB)
      flat  flat%   sum%        cum   cum%
  299.63MB 23.90% 23.90%   299.63MB 23.90%  runtime.malg
   84.51MB  6.74% 30.64%   286.71MB 22.87%  net/http.(*Transport).dialConn
   68.01MB  5.42% 36.06%    67.50MB  5.38%  encoding/json.Marshal
   62.01MB  4.95% 41.01%    86.01MB  6.86%  context.(*cancelCtx).propagateCancel
   60.01MB  4.79% 45.79%    60.01MB  4.79%  internal/sync.runtime_SemacquireMutex
   48.53MB  3.87% 49.67%    48.53MB  3.87%  net.newFD (inline)
   45.50MB  3.63% 53.29%    45.50MB  3.63%  github.com/Okenamay/shorturl.git/internal/audit.(*Auditor).LogEvent
   42.02MB  3.35% 56.65%    42.02MB  3.35%  os.newFile
      37MB  2.95% 59.60%       37MB  2.95%  context.(*cancelCtx).Done
   33.51MB  2.67% 62.27%    33.51MB  2.67%  net/http.Header.Clone (inline)
   29.51MB  2.35% 64.62%    29.51MB  2.35%  net/http.setupRewindBody (inline)
   27.01MB  2.15% 66.78%    27.01MB  2.15%  net/http.send.func1 (inline)
   24.51MB  1.95% 68.73%    24.51MB  1.95%  net/textproto.MIMEHeader.Set (inline)
   23.85MB  1.90% 70.64%    23.85MB  1.90%  runtime.allgadd
   22.79MB  1.82% 72.45%    22.79MB  1.82%  time.newTimer
      21MB  1.68% 74.13%   187.03MB 14.92%  github.com/Okenamay/shorturl.git/internal/audit.(*Auditor).logToFile
   20.50MB  1.64% 75.76%    80.29MB  6.40%  context.WithDeadlineCause
   19.50MB  1.56% 77.32%    39.51MB  3.15%  net/http.NewRequestWithContext
   19.50MB  1.56% 78.87%    19.50MB  1.56%  net/url.parse
      17MB  1.36% 80.23%    35.62MB  2.84%  net/http.(*Transport).getConn
      15MB  1.20% 81.43%    54.60MB  4.35%  net.(*netFD).connect
      13MB  1.04% 82.46%   198.18MB 15.81%  net.(*Dialer).DialContext
      10MB   0.8% 83.26%   302.01MB 24.09%  github.com/Okenamay/shorturl.git/internal/audit.(*Auditor).logToURL
    9.50MB  0.76% 84.02%    34.01MB  2.71%  context.AfterFunc
       9MB  0.72% 84.74%    33.50MB  2.67%  context.withCancel (inline)
    8.50MB  0.68% 85.41%    42.01MB  3.35%  net/http.(*Client).makeHeadersCopier
    8.12MB  0.65% 86.06%     8.12MB  0.65%  net/http.(*wantConnQueue).pushBack (inline)
    7.10MB  0.57% 86.63%     8.60MB  0.69%  net/http.(*Transport).prepareTransportCancel
       7MB  0.56% 87.19%        7MB  0.56%  bytes.NewBuffer (inline)
    6.50MB  0.52% 87.71%     6.50MB  0.52%  syscall.UTF16FromString
    3.50MB  0.28% 87.98%     7.50MB   0.6%  net/http.(*Client).do.func2
    2.50MB   0.2% 88.18%        7MB  0.56%  net.filterAddrList
    1.50MB  0.12% 88.30%   112.73MB  8.99%  net/http.(*Transport).roundTrip
    1.50MB  0.12% 88.42%        6MB  0.48%  context.WithCancel
       1MB  0.08% 88.50%   106.12MB  8.46%  net.(*sysDialer).dialSingle
   -0.50MB  0.04% 88.46%   188.75MB 15.05%  net/http.(*Client).do
         0     0% 88.46%       29MB  2.31%  context.WithCancelCause
         0     0% 88.46%    80.29MB  6.40%  context.WithDeadline (inline)
         0     0% 88.46%    23.24MB  1.85%  context.WithTimeout
         0     0% 88.46%    44.50MB  3.55%  github.com/Okenamay/shorturl.git/internal/audit.BenchmarkLogEventToFile
         0     0% 88.46%    60.01MB  4.79%  internal/sync.(*Mutex).Lock (inline)
         0     0% 88.46%    60.01MB  4.79%  internal/sync.(*Mutex).lockSlow
         0     0% 88.46%       10MB   0.8%  net.(*Resolver).internetAddrList
         0     0% 88.46%       10MB   0.8%  net.(*Resolver).resolveAddrList
         0     0% 88.46%    56.60MB  4.51%  net.(*netFD).dial
         0     0% 88.46%   112.63MB  8.98%  net.(*sysDialer).dialParallel
         0     0% 88.46%   112.63MB  8.98%  net.(*sysDialer).dialSerial
         0     0% 88.46%   105.12MB  8.38%  net.(*sysDialer).dialTCP
         0     0% 88.46%   105.12MB  8.38%  net.(*sysDialer).doDialTCP (inline)
         0     0% 88.46%   105.12MB  8.38%  net.(*sysDialer).doDialTCPProto
         0     0% 88.46%   105.12MB  8.38%  net.internetSocket
         0     0% 88.46%   105.12MB  8.38%  net.socket
         0     0% 88.46%   188.75MB 15.05%  net/http.(*Client).Do (inline)
         0     0% 88.46%   139.74MB 11.15%  net/http.(*Client).send
         0     0% 88.46%   112.73MB  8.99%  net/http.(*Transport).RoundTrip
         0     0% 88.46%   198.18MB 15.81%  net/http.(*Transport).dial
         0     0% 88.46%   286.71MB 22.87%  net/http.(*Transport).dialConnFor
         0     0% 88.46%   286.71MB 22.87%  net/http.(*Transport).startDialConnForLocked.func1
         0     0% 88.46%    24.51MB  1.95%  net/http.Header.Set (inline)
         0     0% 88.46%    33.51MB  2.67%  net/http.cloneOrMakeHeader
         0     0% 88.46%   139.74MB 11.15%  net/http.send
         0     0% 88.46%    19.50MB  1.56%  net/url.Parse
         0     0% 88.46%    48.52MB  3.87%  os.OpenFile
         0     0% 88.46%    48.52MB  3.87%  os.openFileNolog
         0     0% 88.46%   323.47MB 25.80%  runtime.newproc.func1
         0     0% 88.46%   323.47MB 25.80%  runtime.newproc1
         0     0% 88.46%   322.97MB 25.76%  runtime.systemstack
         0     0% 88.46%    60.01MB  4.79%  sync.(*Mutex).Lock (inline)
         0     0% 88.46%     6.50MB  0.52%  syscall.Open
         0     0% 88.46%     6.50MB  0.52%  syscall.UTF16PtrFromString (inline)
         0     0% 88.46%    45.50MB  3.63%  testing.(*B).launch
         0     0% 88.46%    45.50MB  3.63%  testing.(*B).runN
         0     0% 88.46%    22.79MB  1.82%  time.AfterFunc

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_config.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_config.result.pprof
--------------------------------------------------
File: config.test.exe
Build ID: C:\Users\Yamabaka\Dev\shorturl\config.test.exe2025-11-01 22:11:03.7617437 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:31:13 MSK
Showing nodes accounting for 1042.06kB, 50.83% of 2050.06kB total
Dropped 2 nodes (cum <= 10.25kB)
      flat  flat%   sum%        cum   cum%
  528.17kB 25.76% 25.76%   528.17kB 25.76%  regexp.(*bitState).reset
  513.56kB 25.05% 50.81%   513.56kB 25.05%  sync.(*Pool).pinSlow
  512.34kB 24.99% 75.81%   512.34kB 24.99%  crypto/internal/fips140/rsa.map.init.0
 -512.02kB 24.98% 50.83%  -512.02kB 24.98%  crypto/internal/fips140/ecdh.init
         0     0% 50.83%   512.34kB 24.99%  crypto/internal/fips140/rsa.init
         0     0% 50.83%   513.56kB 25.05%  fmt.Fprint
         0     0% 50.83%   513.56kB 25.05%  fmt.Print (inline)
         0     0% 50.83%   513.56kB 25.05%  fmt.newPrinter
         0     0% 50.83%  1041.73kB 50.81%  main.main
         0     0% 50.83%   528.17kB 25.76%  regexp.(*Regexp).MatchString (inline)
         0     0% 50.83%   528.17kB 25.76%  regexp.(*Regexp).backtrack
         0     0% 50.83%   528.17kB 25.76%  regexp.(*Regexp).doExecute
         0     0% 50.83%   528.17kB 25.76%  regexp.(*Regexp).doMatch (inline)
         0     0% 50.83%      513kB 25.02%  runtime.goschedImpl
         0     0% 50.83%      513kB 25.02%  runtime.gosched_m
         0     0% 50.83%  1042.06kB 50.83%  runtime.main
         0     0% 50.83%     -513kB 25.02%  runtime.park_m
         0     0% 50.83%   513.56kB 25.05%  sync.(*Pool).Get
         0     0% 50.83%   513.56kB 25.05%  sync.(*Pool).pin
         0     0% 50.83%  1041.73kB 50.81%  testing.(*M).Run
         0     0% 50.83%   528.17kB 25.76%  testing.newMatcher
         0     0% 50.83%   528.17kB 25.76%  testing.runBenchmarks
         0     0% 50.83%   528.17kB 25.76%  testing.simpleMatch.verify
         0     0% 50.83%   528.17kB 25.76%  testing/internal/testdeps.TestDeps.MatchString

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_logger_zap.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_logger_zap.result.pprof
--------------------------------------------------
File: zap.test.exe
Build ID: C:\Users\Yamabaka\Dev\shorturl\zap.test.exe2025-11-01 22:11:04.3977672 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:31:15 MSK
Showing nodes accounting for -1024.07kB, 36.34% of 2818.33kB total
Dropped 1 node (cum <= 14.09kB)
      flat  flat%   sum%        cum   cum%
 -512.05kB 18.17% 18.17%  -512.05kB 18.17%  runtime.acquireSudog
 -512.02kB 18.17% 36.34%  -512.02kB 18.17%  sync.OnceValue[go.shape.*uint8] (inline)
         0     0% 36.34%  -512.02kB 18.17%  crypto/internal/fips140/ecdsa.init
         0     0% 36.34%  -512.05kB 18.17%  runtime.chanrecv
         0     0% 36.34%  -512.05kB 18.17%  runtime.chanrecv1
         0     0% 36.34%  -512.02kB 18.17%  runtime.doInit
         0     0% 36.34%  -512.02kB 18.17%  runtime.doInit1
         0     0% 36.34%      513kB 18.20%  runtime.gopreempt_m (inline)
         0     0% 36.34%      513kB 18.20%  runtime.goschedImpl
         0     0% 36.34%     -513kB 18.20%  runtime.mcall
         0     0% 36.34%      513kB 18.20%  runtime.morestack
         0     0% 36.34%      513kB 18.20%  runtime.newstack
         0     0% 36.34%     -513kB 18.20%  runtime.park_m
         0     0% 36.34%  -512.05kB 18.17%  runtime.unique_runtime_registerUniqueMapCleanup.func2

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_server_handlers.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_server_handlers.result.pprof
--------------------------------------------------
File: handlers.test.exe
Build ID: C:\Users\Yamabaka\AppData\Local\Temp\go-build3765020729\b001\handlers.test.exe2025-11-01 22:45:35.0571966 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:31:38 MSK
Showing nodes accounting for -387MB, 3.08% of 12584.26MB total
Dropped 95 nodes (cum <= 62.92MB)
      flat  flat%   sum%        cum   cum%
 -300.74MB  2.39%  2.39%  -300.74MB  2.39%  bytes.growSlice
 -222.80MB  1.77%  4.16%  -222.80MB  1.77%  bufio.NewReaderSize (inline)
   64.03MB  0.51%  3.65%    64.03MB  0.51%  encoding/json.(*Decoder).refill
   55.52MB  0.44%  3.21%    55.52MB  0.44%  encoding/json.NewDecoder
   54.52MB  0.43%  2.78%    54.52MB  0.43%  net/http.Header.Clone (inline)
     -53MB  0.42%  3.20%      -53MB  0.42%  crypto/md5.New (inline)
   51.50MB  0.41%  2.79%    51.50MB  0.41%  context.WithValue
   46.01MB  0.37%  2.42%    46.01MB  0.37%  net/http.(*Request).SetPathValue (inline)
   34.51MB  0.27%  2.15%    34.51MB  0.27%  net/textproto.MIMEHeader.Set (inline)
  -34.50MB  0.27%  2.42%   -34.50MB  0.27%  encoding/hex.EncodeToString (inline)
  -30.50MB  0.24%  2.67%   -55.01MB  0.44%  encoding/json.Unmarshal
      20MB  0.16%  2.51%       20MB  0.16%  encoding/json.NewEncoder
      15MB  0.12%  2.39%   -77.50MB  0.62%  github.com/Okenamay/shorturl.git/internal/app/hasher.ShortenURL
      12MB 0.095%  2.29%       12MB 0.095%  net/http/httptest.NewRecorder (inline)
  -11.54MB 0.092%  2.38%   -11.54MB 0.092%  sync.(*Pool).pinSlow
   -9.50MB 0.076%  2.46%    -9.50MB 0.076%  net/url.parse
   -9.50MB 0.075%  2.53%   -11.51MB 0.091%  encoding/json.Marshal
      -8MB 0.064%  2.60%   -13.01MB   0.1%  net/http.readRequest
      -8MB 0.064%  2.66%       -8MB 0.064%  github.com/Okenamay/shorturl.git/internal/app/urlmaker.MakeFullURL (inline)
   -7.50MB  0.06%  2.72%    -7.50MB  0.06%  reflect.growslice
      -6MB 0.048%  2.77%       -6MB 0.048%  bytes.NewReader (inline)
      -6MB 0.048%  2.82%       -6MB 0.048%  encoding/json.(*scanner).pushParseState
   -5.50MB 0.044%  2.86%    -5.50MB 0.044%  crypto/md5.(*digest).Sum (inline)
      -5MB  0.04%  2.90%  -305.74MB  2.43%  bytes.(*Buffer).grow
   -4.50MB 0.036%  2.94%      -47MB  0.37%  github.com/Okenamay/shorturl.git/internal/storage/memselect.ProcessBatchTransaction
       4MB 0.032%  2.90%     2.50MB  0.02%  encoding/json.(*decodeState).object
       4MB 0.032%  2.87%      -11MB 0.087%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkShortenHandler.ShortenHandler.func2
      -4MB 0.032%  2.90%       -4MB 0.032%  net/textproto.(*Reader).ReadLine (inline)
      -4MB 0.032%  2.94%  -220.20MB  1.75%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkJSONHandler.JSONHandler.func3     
   -3.50MB 0.028%  2.96%    -3.50MB 0.028%  io.ReadAll
    3.50MB 0.028%  2.94%  -249.82MB  1.99%  net/http/httptest.NewRequestWithContext
   -3.50MB 0.028%  2.96%    -8.01MB 0.064%  fmt.Sprintf
      -3MB 0.024%  2.99%   -37.51MB   0.3%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkBatchDeleter.BatchDeleter.func1   
      -3MB 0.024%  3.01%   -62.51MB   0.5%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkBatchHandlerTransaction.BatchHandlerTransaction.func1
      -2MB 0.016%  3.03%       -2MB 0.016%  fmt.init.func1
      -2MB 0.016%  3.04%       -2MB 0.016%  unicode/utf16.Encode
   -1.50MB 0.012%  3.06%   -48.52MB  0.39%  net/http.Error
      -1MB 0.0079%  3.06%     3.51MB 0.028%  go.uber.org/zap/internal/stacktrace.Capture
      -1MB 0.0079%  3.07%       -7MB 0.056%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkPingHandler.PingHandler.func2    
   -0.50MB 0.004%  3.08%   -19.02MB  0.15%  github.com/Okenamay/shorturl.git/internal/app/checker.CheckURL
         0     0%  3.08%  -222.80MB  1.77%  bufio.NewReader (inline)
         0     0%  3.08%  -300.24MB  2.39%  bytes.(*Buffer).ReadFrom
         0     0%  3.08%       -8MB 0.064%  bytes.(*Buffer).Write
         0     0%  3.08%     2.50MB  0.02%  bytes.(*Buffer).WriteString
         0     0%  3.08%    77.53MB  0.62%  encoding/json.(*Decoder).Decode
         0     0%  3.08%    64.53MB  0.51%  encoding/json.(*Decoder).readValue
         0     0%  3.08%    17.01MB  0.14%  encoding/json.(*Encoder).Encode
         0     0%  3.08%       -8MB 0.064%  encoding/json.(*decodeState).array
         0     0%  3.08%    -6.50MB 0.052%  encoding/json.(*decodeState).unmarshal
         0     0%  3.08%       -6MB 0.048%  encoding/json.(*decodeState).value
         0     0%  3.08%       -5MB  0.04%  encoding/json.checkValid
         0     0%  3.08%       -6MB 0.048%  encoding/json.stateBeginValue
         0     0%  3.08%    -9.51MB 0.076%  fmt.Fprintln
         0     0%  3.08%    -9.02MB 0.072%  fmt.newPrinter
         0     0%  3.08%   -85.50MB  0.68%  github.com/Okenamay/shorturl.git/internal/app/urlmaker.ProcessURL
         0     0%  3.08%  -166.87MB  1.33%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkBatchDeleter
         0     0%  3.08%  -121.15MB  0.96%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkBatchHandlerTransaction
         0     0%  3.08%      -32MB  0.25%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkJSONHandler.func1
         0     0%  3.08%  -285.52MB  2.27%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkJSONHandler.func2
         0     0%  3.08%       -5MB  0.04%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkPingHandler
         0     0%  3.08%   314.59MB  2.50%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkRedirectHandler
         0     0%  3.08%   145.06MB  1.15%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkRedirectHandler.RedirectHandler.func1
         0     0%  3.08%   -10.01MB  0.08%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkShortenHandler.func1
         0     0%  3.08%   -56.02MB  0.45%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkUserURLsHandler.UserURLsHandler.func3
         0     0%  3.08%    -3.01MB 0.024%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkUserURLsHandler.func1
         0     0%  3.08%   -81.02MB  0.64%  github.com/Okenamay/shorturl.git/internal/server/handlers.BenchmarkUserURLsHandler.func2
         0     0%  3.08%        6MB 0.048%  github.com/Okenamay/shorturl.git/internal/storage/memselect.StorePair
         0     0%  3.08%       -5MB  0.04%  github.com/Okenamay/shorturl.git/internal/storage/memselect.StorePairTransaction
         0     0%  3.08%  -111.67MB  0.89%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0%  3.08%  -203.17MB  1.61%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0%  3.08%    46.01MB  0.37%  github.com/go-chi/chi/v5.setPathValue
         0     0%  3.08%     3.51MB 0.028%  go.uber.org/zap.(*Logger).Check (inline)
         0     0%  3.08%     3.51MB 0.028%  go.uber.org/zap.(*Logger).check
         0     0%  3.08%    -9.02MB 0.072%  go.uber.org/zap.(*SugaredLogger).Infof (partial-inline)
         0     0%  3.08%    -8.51MB 0.068%  go.uber.org/zap.(*SugaredLogger).log
         0     0%  3.08%    -8.01MB 0.064%  go.uber.org/zap.getMessage
         0     0%  3.08%     2.50MB  0.02%  go.uber.org/zap/internal/pool.(*Pool[go.shape.*uint8]).Get (inline)
         0     0%  3.08%    -4.01MB 0.032%  go.uber.org/zap/zapcore.(*CheckedEntry).Write
         0     0%  3.08%    -4.01MB 0.032%  go.uber.org/zap/zapcore.(*ioCore).Write
         0     0%  3.08%    -2.51MB  0.02%  go.uber.org/zap/zapcore.(*jsonEncoder).EncodeEntry
         0     0%  3.08%    -2.01MB 0.016%  go.uber.org/zap/zapcore.(*jsonEncoder).clone
         0     0%  3.08%       -2MB 0.016%  go.uber.org/zap/zapcore.(*lockedWriteSyncer).Write
         0     0%  3.08%       -2MB 0.016%  internal/poll.(*FD).Write
         0     0%  3.08%       -2MB 0.016%  internal/poll.(*FD).writeConsole
         0     0%  3.08%  -203.17MB  1.61%  net/http.HandlerFunc.ServeHTTP
         0     0%  3.08%    34.51MB  0.27%  net/http.Header.Set (partial-inline)
         0     0%  3.08%   -13.01MB   0.1%  net/http.ReadRequest
         0     0%  3.08%    -3.51MB 0.028%  net/http.newTextprotoReader
         0     0%  3.08%       -8MB 0.064%  net/http/httptest.(*ResponseRecorder).Write
         0     0%  3.08%    54.52MB  0.43%  net/http/httptest.(*ResponseRecorder).WriteHeader
         0     0%  3.08%     2.50MB  0.02%  net/http/httptest.(*ResponseRecorder).WriteString
         0     0%  3.08%  -249.82MB  1.99%  net/http/httptest.NewRequest (inline)
         0     0%  3.08%    -9.50MB 0.076%  net/url.ParseRequestURI
         0     0%  3.08%       -2MB 0.016%  os.(*File).Write
         0     0%  3.08%       -2MB 0.016%  os.(*File).write (inline)
         0     0%  3.08%    -7.50MB  0.06%  reflect.Value.Grow
         0     0%  3.08%    -7.50MB  0.06%  reflect.Value.grow
         0     0%  3.08%    -7.52MB  0.06%  sync.(*Pool).Get
         0     0%  3.08%    -5.02MB  0.04%  sync.(*Pool).Put
         0     0%  3.08%   -11.54MB 0.092%  sync.(*Pool).pin
         0     0%  3.08%  -389.50MB  3.10%  testing.(*B).launch
         0     0%  3.08%     -390MB  3.10%  testing.(*B).runN

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_server_router.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_server_router.result.pprof
--------------------------------------------------
File: router.test.exe
Build ID: C:\Users\Yamabaka\AppData\Local\Temp\go-build3937142972\b001\router.test.exe2025-11-01 22:45:57.9555069 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:31:45 MSK
Showing nodes accounting for 204.90MB, 3.28% of 6254.29MB total
Dropped 47 nodes (cum <= 31.27MB)
      flat  flat%   sum%        cum   cum%
  101.90MB  1.63%  1.63%   101.90MB  1.63%  bufio.NewReaderSize (inline)
   27.01MB  0.43%  2.06%    27.01MB  0.43%  net/textproto.MIMEHeader.Add (inline)
   25.01MB   0.4%  2.46%    25.01MB   0.4%  net/http.(*Request).WithContext (partial-inline)
   18.50MB   0.3%  2.76%    18.50MB   0.3%  encoding/base64.(*Encoding).EncodeToString
   14.50MB  0.23%  2.99%    14.50MB  0.23%  crypto/internal/fips140/sha256.New (inline)
  -13.50MB  0.22%  2.77%   -13.50MB  0.22%  github.com/golang-jwt/jwt/v4.NewWithClaims (inline)
  -11.50MB  0.18%  2.59%   -11.50MB  0.18%  strings.(*Builder).grow
      11MB  0.18%  2.76%    -2.51MB  0.04%  net/http.readRequest
   10.50MB  0.17%  2.93%    29.50MB  0.47%  github.com/golang-jwt/jwt/v4.(*SigningMethodHMAC).Sign
   -8.50MB  0.14%  2.80%    -8.50MB  0.14%  net/url.parse
   -8.50MB  0.14%  2.66%       -9MB  0.14%  encoding/json.mapEncoder.encode
    8.50MB  0.14%  2.80%    10.50MB  0.17%  github.com/golang-jwt/jwt/v4.NumericDate.MarshalJSON
   -8.02MB  0.13%  2.67%    -8.02MB  0.13%  sync.(*Pool).pinSlow
       8MB  0.13%  2.80%        8MB  0.13%  net/http.Header.Clone (inline)
    7.50MB  0.12%  2.92%    18.48MB   0.3%  github.com/Okenamay/shorturl.git/internal/app/middleware/auth.buildJWTString
       7MB  0.11%  3.03%        7MB  0.11%  context.WithValue
    6.50MB   0.1%  3.13%     6.50MB   0.1%  github.com/golang-jwt/jwt/v4.NewNumericDate (inline)
       6MB 0.096%  3.23%        6MB 0.096%  crypto/internal/fips140/sha256.(*Digest).Sum
    5.50MB 0.088%  3.32%     5.50MB 0.088%  strconv.FormatFloat (inline)
    4.50MB 0.072%  3.39%     4.50MB 0.072%  strings.NewReader (inline)
   -3.50MB 0.056%  3.33%    -3.50MB 0.056%  strconv.formatBits
      -3MB 0.048%  3.28%   111.40MB  1.78%  net/http/httptest.NewRequestWithContext
      -3MB 0.048%  3.24%       -3MB 0.048%  net/textproto.MIMEHeader.Set (inline)
    2.50MB  0.04%  3.28%       67MB  1.07%  github.com/Okenamay/shorturl.git/internal/server/router.NewRouter.WithLogging.func1.1
    2.50MB  0.04%  3.32%       -2MB 0.032%  net/http.Error
      -2MB 0.032%  3.28%       -2MB 0.032%  runtime.allocm
       2MB 0.032%  3.32%     2.50MB  0.04%  bytes.(*Buffer).grow
   -1.50MB 0.024%  3.29%    -1.50MB 0.024%  github.com/google/uuid.UUID.String
    1.50MB 0.024%  3.32%     1.50MB 0.024%  github.com/google/uuid.NewRandomFromReader
      -1MB 0.016%  3.30%       -1MB 0.016%  fmt.init.func1
       1MB 0.016%  3.32%        1MB 0.016%  encoding/json.init.func1
      -1MB 0.016%  3.30%       -1MB 0.016%  reflect.copyVal
      -1MB 0.016%  3.28%       -1MB 0.016%  net/textproto.(*Reader).ReadLine (inline)
   -0.50MB 0.008%  3.28%    64.50MB  1.03%  github.com/Okenamay/shorturl.git/internal/server/router.NewRouter.Authenticator.func2.1
         0     0%  3.28%   101.90MB  1.63%  bufio.NewReader (inline)
         0     0%  3.28%        2MB 0.032%  bytes.(*Buffer).Write
         0     0%  3.28%    14.50MB  0.23%  crypto.Hash.New
         0     0%  3.28%    14.50MB  0.23%  crypto/hmac.New.UnwrapNew[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }].func1
         0     0%  3.28%        6MB 0.096%  crypto/internal/fips140/hmac.(*HMAC).Sum
         0     0%  3.28%    14.50MB  0.23%  crypto/sha256.New
         0     0%  3.28%     3.50MB 0.056%  encoding/json.(*encodeState).marshal
         0     0%  3.28%     3.50MB 0.056%  encoding/json.(*encodeState).reflectValue
         0     0%  3.28%     1.50MB 0.024%  encoding/json.appendCompact
         0     0%  3.28%       12MB  0.19%  encoding/json.marshalerEncoder
         0     0%  3.28%    -6.52MB   0.1%  encoding/json.newEncodeState
         0     0%  3.28%     1.50MB 0.024%  encoding/json.newScanner
         0     0%  3.28%    12.50MB   0.2%  encoding/json.structEncoder.encode
         0     0%  3.28%     1.50MB 0.024%  fmt.Fprintln
         0     0%  3.28%     1.50MB 0.024%  github.com/Okenamay/shorturl.git/internal/logger/zap.(*loggingResponseWriter).Write
         0     0%  3.28%        8MB  0.13%  github.com/Okenamay/shorturl.git/internal/logger/zap.(*loggingResponseWriter).WriteHeader
         0     0%  3.28%     1.02MB 0.016%  github.com/Okenamay/shorturl.git/internal/server/router.BenchmarkRoutes
         0     0%  3.28%   194.40MB  3.11%  github.com/Okenamay/shorturl.git/internal/server/router.BenchmarkRoutes.func1
         0     0%  3.28%     3.50MB 0.056%  github.com/Okenamay/shorturl.git/internal/server/router.NewRouter.PingHandler.func3
         0     0%  3.28%       12MB  0.19%  github.com/Okenamay/shorturl.git/internal/server/router.NewRouter.UserURLsHandler.func6
         0     0%  3.28%       83MB  1.33%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0%  3.28%    10.01MB  0.16%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0%  3.28%    17.98MB  0.29%  github.com/golang-jwt/jwt/v4.(*Token).SignedString
         0     0%  3.28%   -11.02MB  0.18%  github.com/golang-jwt/jwt/v4.(*Token).SigningString
         0     0%  3.28%    18.50MB   0.3%  github.com/golang-jwt/jwt/v4.EncodeSegment (inline)
         0     0%  3.28%     1.50MB 0.024%  github.com/google/uuid.New
         0     0%  3.28%     1.50MB 0.024%  github.com/google/uuid.NewRandom
         0     0%  3.28%        6MB 0.096%  net/http.(*Cookie).String
         0     0%  3.28%       67MB  1.07%  net/http.HandlerFunc.ServeHTTP
         0     0%  3.28%    27.01MB  0.43%  net/http.Header.Add (partial-inline)
         0     0%  3.28%       -3MB 0.048%  net/http.Header.Set (inline)
         0     0%  3.28%    -5.50MB 0.088%  net/http.NotFound
         0     0%  3.28%    -2.51MB  0.04%  net/http.ReadRequest
         0     0%  3.28%    32.01MB  0.51%  net/http.SetCookie
         0     0%  3.28%    -3.01MB 0.048%  net/http.newTextprotoReader
         0     0%  3.28%     1.50MB 0.024%  net/http/httptest.(*ResponseRecorder).Write
         0     0%  3.28%        8MB  0.13%  net/http/httptest.(*ResponseRecorder).WriteHeader
         0     0%  3.28%   111.40MB  1.78%  net/http/httptest.NewRequest (inline)
         0     0%  3.28%    -8.50MB  0.14%  net/url.ParseRequestURI
         0     0%  3.28%        4MB 0.064%  reflect.(*MapIter).Key
         0     0%  3.28%       -5MB  0.08%  reflect.(*MapIter).Value
         0     0%  3.28%    -1.50MB 0.024%  runtime.mcall
         0     0%  3.28%       -2MB 0.032%  runtime.newm
         0     0%  3.28%    -1.50MB 0.024%  runtime.park_m
         0     0%  3.28%    -1.50MB 0.024%  runtime.resetspinning
         0     0%  3.28%    -1.50MB 0.024%  runtime.schedule
         0     0%  3.28%       -2MB 0.032%  runtime.startm
         0     0%  3.28%       -2MB 0.032%  runtime.wakep
         0     0%  3.28%    -3.50MB 0.056%  strconv.FormatInt
         0     0%  3.28%   -11.50MB  0.18%  strings.(*Builder).Grow
         0     0%  3.28%   -17.50MB  0.28%  strings.Join
         0     0%  3.28%    -7.52MB  0.12%  sync.(*Pool).Get
         0     0%  3.28%    -8.02MB  0.13%  sync.(*Pool).pin
         0     0%  3.28%     1.02MB 0.016%  testing.(*B).Run
         0     0%  3.28%   194.40MB  3.11%  testing.(*B).launch
         0     0%  3.28%     1.02MB 0.016%  testing.(*B).run1.func1
         0     0%  3.28%   195.42MB  3.12%  testing.(*B).runN

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_storage_memselect.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_storage_memselect.result.pprof
--------------------------------------------------
File: memselect.test.exe
Build ID: C:\Users\Yamabaka\AppData\Local\Temp\go-build2006939613\b001\memselect.test.exe2025-11-01 22:46:05.5403801 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:32:12 MSK
Showing nodes accounting for -76.51MB, 19.39% of 394.60MB total
Dropped 10 nodes (cum <= 1.97MB)
      flat  flat%   sum%        cum   cum%
  -30.50MB  7.73%  7.73%   -30.50MB  7.73%  crypto/md5.New (inline)
  -30.01MB  7.60% 15.33%   -30.01MB  7.60%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.(*URLMap).Set
     -26MB  6.59% 21.92%      -26MB  6.59%  encoding/hex.EncodeToString (inline)
      14MB  3.55% 18.38%   -49.50MB 12.55%  github.com/Okenamay/shorturl.git/internal/app/hasher.ShortenURL
       7MB  1.77% 16.60%   -48.50MB 12.29%  github.com/Okenamay/shorturl.git/internal/storage/memselect.ProcessBatchTransaction
      -7MB  1.77% 18.38%       -7MB  1.77%  crypto/md5.(*digest).Sum (inline)
      -6MB  1.52% 19.90%       -6MB  1.52%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.NewURLMap (inline)
   -3.50MB  0.89% 20.78%    -3.50MB  0.89%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.(*URLMap).GetAll
    3.50MB  0.89% 19.90%     3.50MB  0.89%  github.com/Okenamay/shorturl.git/internal/app/urlmaker.MakeFullURL (inline)
    2.50MB  0.63% 19.26%     2.50MB  0.63%  unicode/utf16.Encode
       2MB  0.51% 18.76%     1.50MB  0.38%  go.uber.org/zap/internal/stacktrace.Capture
   -1.50MB  0.38% 19.14%    -1.50MB  0.38%  encoding/json.Marshal
      -1MB  0.25% 19.39%       -1MB  0.25%  sync.(*Pool).pinSlow
       1MB  0.25% 19.14%        1MB  0.25%  syscall.UTF16FromString
   -0.50MB  0.13% 19.26%    -0.50MB  0.13%  regexp/syntax.(*compiler).inst (inline)
   -0.50MB  0.13% 19.39%    -0.50MB  0.13%  net/http.init
    0.50MB  0.13% 19.26%     0.50MB  0.13%  os.newFileStatFromWin32FileAttributeData (inline)
    0.50MB  0.13% 19.14%     0.50MB  0.13%  runtime.acquireSudog
   -0.50MB  0.13% 19.26%    -4.50MB  1.14%  github.com/Okenamay/shorturl.git/internal/storage/savefile.SaveFile
   -0.50MB  0.13% 19.39%    -0.50MB  0.13%  internal/filepathlite.Dir
    0.50MB  0.13% 19.26%     0.50MB  0.13%  go.uber.org/zap/buffer.(*Buffer).String (inline)
   -0.50MB  0.13% 19.39%    -0.50MB  0.13%  os.newFile
         0     0% 19.39%      -46MB 11.66%  github.com/Okenamay/shorturl.git/internal/app/urlmaker.ProcessURL
         0     0% 19.39%      -46MB 11.66%  github.com/Okenamay/shorturl.git/internal/storage/memselect.BenchmarkProcessBatchTransaction.func1   
         0     0% 19.39%    -3.50MB  0.89%  github.com/Okenamay/shorturl.git/internal/storage/memselect.BenchmarkProcessBatchTransaction.func2   
         0     0% 19.39%   -21.51MB  5.45%  github.com/Okenamay/shorturl.git/internal/storage/memselect.BenchmarkStorePair.func1
         0     0% 19.39%       -5MB  1.27%  github.com/Okenamay/shorturl.git/internal/storage/memselect.BenchmarkStorePair.func2
         0     0% 19.39%   -21.51MB  5.45%  github.com/Okenamay/shorturl.git/internal/storage/memselect.StorePair
         0     0% 19.39%    -8.50MB  2.15%  github.com/Okenamay/shorturl.git/internal/storage/memselect.StorePairTransaction
         0     0% 19.39%        1MB  0.25%  go.uber.org/zap.(*Logger).Check (inline)
         0     0% 19.39%        1MB  0.25%  go.uber.org/zap.(*Logger).check
         0     0% 19.39%        4MB  1.01%  go.uber.org/zap.(*SugaredLogger).Info (inline)
         0     0% 19.39%        4MB  1.01%  go.uber.org/zap.(*SugaredLogger).log
         0     0% 19.39%     0.50MB  0.13%  go.uber.org/zap/buffer.(*Buffer).Free (inline)
         0     0% 19.39%       -1MB  0.25%  go.uber.org/zap/buffer.Pool.Get
         0     0% 19.39%     0.50MB  0.13%  go.uber.org/zap/buffer.Pool.put (inline)
         0     0% 19.39%    -1.50MB  0.38%  go.uber.org/zap/internal/pool.(*Pool[go.shape.*uint8]).Get (inline)
         0     0% 19.39%     0.50MB  0.13%  go.uber.org/zap/internal/pool.(*Pool[go.shape.*uint8]).Put (inline)
         0     0% 19.39%    -0.50MB  0.13%  go.uber.org/zap/zapcore.(*CheckedEntry).AddCore (inline)
         0     0% 19.39%        3MB  0.76%  go.uber.org/zap/zapcore.(*CheckedEntry).Write
         0     0% 19.39%    -0.50MB  0.13%  go.uber.org/zap/zapcore.(*ioCore).Check
         0     0% 19.39%        3MB  0.76%  go.uber.org/zap/zapcore.(*ioCore).Write
         0     0% 19.39%    -0.50MB  0.13%  go.uber.org/zap/zapcore.(*jsonEncoder).clone
         0     0% 19.39%     2.50MB  0.63%  go.uber.org/zap/zapcore.(*lockedWriteSyncer).Write
         0     0% 19.39%    -0.50MB  0.13%  go.uber.org/zap/zapcore.(*sampler).Check
         0     0% 19.39%     0.50MB  0.13%  go.uber.org/zap/zapcore.EntryCaller.TrimmedPath
         0     0% 19.39%     0.50MB  0.13%  go.uber.org/zap/zapcore.ShortCallerEncoder
         0     0% 19.39%    -0.50MB  0.13%  go.uber.org/zap/zapcore.getCheckedEntry
         0     0% 19.39%    -0.50MB  0.13%  gopkg.in/yaml%2ev3.init
         0     0% 19.39%     2.50MB  0.63%  internal/poll.(*FD).Write
         0     0% 19.39%     2.50MB  0.63%  internal/poll.(*FD).writeConsole
         0     0% 19.39%     2.50MB  0.63%  os.(*File).Write
         0     0% 19.39%     2.50MB  0.63%  os.(*File).write (inline)
         0     0% 19.39%        1MB  0.25%  os.MkdirAll
         0     0% 19.39%        1MB  0.25%  os.Stat
         0     0% 19.39%        1MB  0.25%  os.stat
         0     0% 19.39%        1MB  0.25%  os.statNolog (inline)
         0     0% 19.39%    -0.50MB  0.13%  path/filepath.Dir (inline)
         0     0% 19.39%    -0.50MB  0.13%  regexp.Compile (inline)
         0     0% 19.39%    -0.50MB  0.13%  regexp.MustCompile
         0     0% 19.39%    -0.50MB  0.13%  regexp.compile
         0     0% 19.39%    -0.50MB  0.13%  regexp/syntax.(*compiler).cap (inline)
         0     0% 19.39%    -0.50MB  0.13%  regexp/syntax.(*compiler).compile
         0     0% 19.39%    -0.50MB  0.13%  regexp/syntax.Compile
         0     0% 19.39%       -1MB  0.25%  runtime.doInit (inline)
         0     0% 19.39%       -1MB  0.25%  runtime.doInit1
         0     0% 19.39%     0.50MB  0.13%  runtime.gcBgMarkWorker
         0     0% 19.39%     0.50MB  0.13%  runtime.gcMarkDone
         0     0% 19.39%     0.50MB  0.13%  runtime.goschedImpl
         0     0% 19.39%     0.50MB  0.13%  runtime.gosched_m
         0     0% 19.39%       -1MB  0.25%  runtime.main
         0     0% 19.39%    -0.50MB  0.13%  runtime.park_m
         0     0% 19.39%     0.50MB  0.13%  runtime.semacquire (inline)
         0     0% 19.39%     0.50MB  0.13%  runtime.semacquire1
         0     0% 19.39%    -1.50MB  0.38%  sync.(*Pool).Get
         0     0% 19.39%     0.50MB  0.13%  sync.(*Pool).Put
         0     0% 19.39%       -1MB  0.25%  sync.(*Pool).pin
         0     0% 19.39%     0.50MB  0.13%  syscall.Open
         0     0% 19.39%        1MB  0.25%  syscall.UTF16PtrFromString (inline)
         0     0% 19.39%   -75.51MB 19.14%  testing.(*B).launch
         0     0% 19.39%    -0.50MB  0.13%  testing.(*B).run1.func1
         0     0% 19.39%   -76.01MB 19.26%  testing.(*B).runN

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_storage_memstorage.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_storage_memstorage.result.pprof
--------------------------------------------------
File: memstorage.test.exe
Build ID: C:\Users\Yamabaka\Dev\shorturl\memstorage.test.exe2025-11-01 22:12:07.5185811 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:32:21 MSK
Showing nodes accounting for -1224.44MB, 33.94% of 3608.13MB total
Dropped 50 nodes (cum <= 18.04MB)
      flat  flat%   sum%        cum   cum%
-1251.69MB 34.69% 34.69% -1251.69MB 34.69%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.(*URLMap).GetAll
   46.25MB  1.28% 33.41%    46.25MB  1.28%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.(*URLMap).Set
  -25.50MB  0.71% 34.12%   -25.50MB  0.71%  fmt.Sprintf
   10.50MB  0.29% 33.82%    58.69MB  1.63%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.BenchmarkSet
   -3.50MB 0.097% 33.92%   -30.44MB  0.84%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.BenchmarkConcurrentSetGet.func1
   -0.50MB 0.014% 33.94%    -1194MB 33.09%  testing.(*B).runN
         0     0% 33.94% -1251.69MB 34.69%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.BenchmarkGetAll
         0     0% 33.94%   -30.44MB  0.84%  testing.(*B).RunParallel.func1
         0     0% 33.94% -1193.50MB 33.08%  testing.(*B).launch

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_storage_migration.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_storage_migration.result.pprof
--------------------------------------------------
File: migration.test.exe
Build ID: C:\Users\Yamabaka\AppData\Local\Temp\go-build2873196207\b001\migration.test.exe2025-11-01 22:46:40.499207 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:32:30 MSK
Showing nodes accounting for -250.62MB, 98.62% of 254.13MB total
Dropped 14 nodes (cum <= 1.27MB)
      flat  flat%   sum%        cum   cum%
 -141.28MB 55.59% 55.59%  -141.28MB 55.59%  regexp/syntax.(*compiler).inst (inline)
  -25.01MB  9.84% 65.43%   -25.01MB  9.84%  regexp/syntax.(*parser).maybeConcat
     -14MB  5.51% 70.94%   -19.08MB  7.51%  regexp.(*Regexp).replaceAll
     -10MB  3.94% 74.88%      -10MB  3.94%  regexp/syntax.(*parser).newRegexp (inline)
   -8.50MB  3.35% 78.22%    -8.50MB  3.35%  unicode/utf8.AppendRune (inline)
   -8.02MB  3.15% 81.38%    -8.02MB  3.15%  runtime.allocm
      -7MB  2.75% 84.13%       -7MB  2.75%  regexp.QuoteMeta
      -6MB  2.36% 86.50%       -9MB  3.54%  time.NewTimer
   -5.70MB  2.24% 88.74%    -5.70MB  2.24%  github.com/pashagolub/pgxmock/v3.(*pgxmock).ExpectExec
   -4.11MB  1.62% 90.36%    -4.11MB  1.62%  regexp.(*bitState).reset
      -4MB  1.57% 91.93%  -194.29MB 76.45%  regexp.compile
      -4MB  1.57% 93.51%   -39.01MB 15.35%  regexp/syntax.parse
   -3.50MB  1.38% 94.88%   -22.58MB  8.89%  regexp.(*Regexp).ReplaceAllString
      -3MB  1.18% 96.06%       -3MB  1.18%  time.newTimer
      -2MB  0.79% 96.85%  -228.90MB 90.08%  github.com/Okenamay/shorturl.git/internal/storage/migration.MigrateLauncher
   -1.50MB  0.59% 97.44%       -2MB  0.79%  fmt.Sprintf
   -1.50MB  0.59% 98.03%    -1.50MB  0.59%  regexp.(*Regexp).expand
      -1MB  0.39% 98.42%    -2.50MB  0.98%  regexp/syntax.(*compiler).init (inline)
   -0.50MB   0.2% 98.62%    -2.50MB  0.98%  github.com/pashagolub/pgxmock/v3.NewResult (inline)
         0     0% 98.62%  -225.59MB 88.77%  github.com/Okenamay/shorturl.git/internal/storage/migration.BenchmarkMigrateLauncher.func1
         0     0% 98.62%   -18.52MB  7.29%  github.com/Okenamay/shorturl.git/internal/storage/migration.BenchmarkMigrateLauncher.func2
         0     0% 98.62%       -9MB  3.54%  github.com/pashagolub/pgxmock/v3.(*commonExpectation).waitForDelay
         0     0% 98.62%  -226.90MB 89.29%  github.com/pashagolub/pgxmock/v3.(*pgxmock).Exec
         0     0% 98.62%  -217.90MB 85.75%  github.com/pashagolub/pgxmock/v3.(*pgxmock).Exec.func1
         0     0% 98.62%  -217.90MB 85.75%  github.com/pashagolub/pgxmock/v3.QueryMatcherFunc.Match
         0     0% 98.62%  -217.90MB 85.75%  github.com/pashagolub/pgxmock/v3.findExpectationFunc[go.shape.*github.com/pashagolub/pgxmock/v3.ExpectedExec,go.shape.struct { github.com/pashagolub/pgxmock/v3.commonExpectation; github.com/pashagolub/pgxmock/v3.queryBasedExpectation; github.com/pashagolub/pgxmock/v3.result github.com/jackc/pgx/v5/pgconn.CommandTag }]
         0     0% 98.62%  -217.90MB 85.75%  github.com/pashagolub/pgxmock/v3.init.func1
         0     0% 98.62%   -22.58MB  8.89%  github.com/pashagolub/pgxmock/v3.stripQuery
         0     0% 98.62%    -1.50MB  0.59%  regexp.(*Regexp).ReplaceAllString.func1
         0     0% 98.62%    -4.61MB  1.82%  regexp.(*Regexp).backtrack
         0     0% 98.62%    -4.61MB  1.82%  regexp.(*Regexp).doExecute
         0     0% 98.62%  -194.29MB 76.45%  regexp.Compile (inline)
         0     0% 98.62%    -8.50MB  3.35%  regexp/syntax.(*Prog).Prefix
         0     0% 98.62%  -139.78MB 55.00%  regexp/syntax.(*compiler).compile
         0     0% 98.62%  -139.78MB 55.00%  regexp/syntax.(*compiler).rune
         0     0% 98.62%    -1.50MB  0.59%  regexp/syntax.(*parser).concat
         0     0% 98.62%   -33.51MB 13.19%  regexp/syntax.(*parser).literal
         0     0% 98.62%   -23.51MB  9.25%  regexp/syntax.(*parser).push
         0     0% 98.62%  -142.28MB 55.99%  regexp/syntax.Compile
         0     0% 98.62%   -39.01MB 15.35%  regexp/syntax.Parse (inline)
         0     0% 98.62%    -8.02MB  3.15%  runtime.mcall
         0     0% 98.62%    -8.02MB  3.15%  runtime.newm
         0     0% 98.62%    -8.02MB  3.15%  runtime.park_m
         0     0% 98.62%    -8.02MB  3.15%  runtime.resetspinning
         0     0% 98.62%    -8.02MB  3.15%  runtime.schedule
         0     0% 98.62%    -8.02MB  3.15%  runtime.startm
         0     0% 98.62%    -8.02MB  3.15%  runtime.wakep
         0     0% 98.62%    -8.50MB  3.35%  strings.(*Builder).WriteRune
         0     0% 98.62%  -244.11MB 96.06%  testing.(*B).launch
         0     0% 98.62%  -244.11MB 96.06%  testing.(*B).runN
         0     0% 98.62%       -9MB  3.54%  time.After (inline)

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_storage_savefile.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_storage_savefile.result.pprof
--------------------------------------------------
File: savefile.test.exe
Build ID: C:\Users\Yamabaka\AppData\Local\Temp\go-build193042529\b001\savefile.test.exe2025-11-01 22:46:41.7683208 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:32:34 MSK
Showing nodes accounting for -79.32MB, 18.63% of 425.77MB total
Dropped 1 node (cum <= 2.13MB)
      flat  flat%   sum%        cum   cum%
  -56.94MB 13.37% 13.37%   -57.94MB 13.61%  os.readFileContents
   32.60MB  7.66%  5.72%   -46.41MB 10.90%  github.com/Okenamay/shorturl.git/internal/storage/savefile.LoadFile
  -19.55MB  4.59% 10.31%   -19.55MB  4.59%  bytes.genSplit
  -19.50MB  4.58% 14.89%   -33.99MB  7.98%  github.com/Okenamay/shorturl.git/internal/storage/savefile.SaveFile
     -17MB  3.99% 18.88%   -17.50MB  4.11%  encoding/json.Marshal
      -7MB  1.64% 20.53%       -3MB   0.7%  encoding/json.(*decodeState).object
   -6.02MB  1.41% 21.94%    -6.02MB  1.41%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.(*URLMap).Set
       4MB  0.94% 21.00%        4MB  0.94%  encoding/json.(*decodeState).literalStore
    3.50MB  0.82% 20.18%        2MB  0.47%  encoding/json.Unmarshal
    3.01MB  0.71% 19.47%     3.01MB  0.71%  github.com/Okenamay/shorturl.git/internal/storage/memstorage.(*URLMap).GetAll
    2.50MB  0.59% 18.88%     2.50MB  0.59%  os.newFile
    1.60MB  0.37% 18.51%     1.60MB  0.37%  os.init.func1
    1.50MB  0.35% 18.16%     1.50MB  0.35%  encoding/json.(*scanner).pushParseState
      -1MB  0.23% 18.39%       -1MB  0.23%  os.newFileStatFromGetFileInformationByHandle
   -0.52MB  0.12% 18.51%    -0.52MB  0.12%  regexp.(*bitState).reset
   -0.50MB  0.12% 18.63%    -0.50MB  0.12%  sync.(*Pool).pinSlow
   -0.50MB  0.12% 18.75%    -0.50MB  0.12%  runtime.allocm
    0.50MB  0.12% 18.63%     0.50MB  0.12%  syscall.UTF16FromString
   -0.50MB  0.12% 18.75%    -0.50MB  0.12%  internal/filepathlite.Dir
    0.50MB  0.12% 18.63%     0.50MB  0.12%  runtime.gcBgMarkWorker
         0     0% 18.63%   -19.55MB  4.59%  bytes.Split (inline)
         0     0% 18.63%       -3MB   0.7%  encoding/json.(*decodeState).unmarshal
         0     0% 18.63%       -3MB   0.7%  encoding/json.(*decodeState).value
         0     0% 18.63%     1.50MB  0.35%  encoding/json.checkValid
         0     0% 18.63%    -0.50MB  0.12%  encoding/json.newEncodeState
         0     0% 18.63%     1.50MB  0.35%  encoding/json.stateBeginValue
         0     0% 18.63%   -46.91MB 11.02%  github.com/Okenamay/shorturl.git/internal/storage/savefile.BenchmarkLoadFile
         0     0% 18.63%   -33.49MB  7.87%  github.com/Okenamay/shorturl.git/internal/storage/savefile.BenchmarkSaveFile
         0     0% 18.63%    -0.50MB  0.12%  github.com/Okenamay/shorturl.git/internal/storage/savefile.prepareBenchFile
         0     0% 18.63%    -0.52MB  0.12%  main.main
         0     0% 18.63%     1.60MB  0.37%  os.(*File).Readdirnames
         0     0% 18.63%       -1MB  0.23%  os.(*File).Stat
         0     0% 18.63%     1.60MB  0.37%  os.(*File).readdir
         0     0% 18.63%     2.50MB  0.59%  os.Open (inline)
         0     0% 18.63%        3MB   0.7%  os.OpenFile
         0     0% 18.63%   -62.44MB 14.67%  os.ReadFile
         0     0% 18.63%     1.60MB  0.37%  os.RemoveAll (inline)
         0     0% 18.63%        3MB   0.7%  os.openFileNolog
         0     0% 18.63%     1.60MB  0.37%  os.removeAll
         0     0% 18.63%       -1MB  0.23%  os.statHandle
         0     0% 18.63%    -0.50MB  0.12%  path/filepath.Dir (inline)
         0     0% 18.63%    -0.52MB  0.12%  regexp.(*Regexp).MatchString (inline)
         0     0% 18.63%    -0.52MB  0.12%  regexp.(*Regexp).backtrack
         0     0% 18.63%    -0.52MB  0.12%  regexp.(*Regexp).doExecute
         0     0% 18.63%    -0.52MB  0.12%  regexp.(*Regexp).doMatch (inline)
         0     0% 18.63%    -0.52MB  0.12%  runtime.main
         0     0% 18.63%    -0.50MB  0.12%  runtime.mcall
         0     0% 18.63%    -0.50MB  0.12%  runtime.newm
         0     0% 18.63%    -0.50MB  0.12%  runtime.park_m
         0     0% 18.63%    -0.50MB  0.12%  runtime.resetspinning
         0     0% 18.63%    -0.50MB  0.12%  runtime.schedule
         0     0% 18.63%    -0.50MB  0.12%  runtime.startm
         0     0% 18.63%    -0.50MB  0.12%  runtime.wakep
         0     0% 18.63%     1.09MB  0.26%  sync.(*Pool).Get
         0     0% 18.63%    -0.50MB  0.12%  sync.(*Pool).pin
         0     0% 18.63%     0.50MB  0.12%  syscall.Open
         0     0% 18.63%     0.50MB  0.12%  syscall.UTF16PtrFromString (inline)
         0     0% 18.63%    -0.52MB  0.12%  testing.(*B).Run
         0     0% 18.63%   -78.83MB 18.52%  testing.(*B).launch
         0     0% 18.63%   -79.32MB 18.63%  testing.(*B).runN
         0     0% 18.63%     1.60MB  0.37%  testing.(*B).runN.func1
         0     0% 18.63%    -0.52MB  0.12%  testing.(*M).Run
         0     0% 18.63%     1.60MB  0.37%  testing.(*common).Cleanup.func1
         0     0% 18.63%     1.60MB  0.37%  testing.(*common).TempDir.func2
         0     0% 18.63%     1.60MB  0.37%  testing.(*common).runCleanup
         0     0% 18.63%    -0.52MB  0.12%  testing.(*matcher).fullName
         0     0% 18.63%     1.60MB  0.37%  testing.removeAll
         0     0% 18.63%    -0.52MB  0.12%  testing.runBenchmarks
         0     0% 18.63%    -0.52MB  0.12%  testing.runBenchmarks.func1
         0     0% 18.63%    -0.52MB  0.12%  testing.simpleMatch.matches
         0     0% 18.63%    -0.52MB  0.12%  testing/internal/testdeps.TestDeps.MatchString

==================================================
Processing pair:
  Base:   ./github.com_Okenamay_shorturl.git_internal_worker.base.pprof
  Result: ./github.com_Okenamay_shorturl.git_internal_worker.result.pprof
--------------------------------------------------
File: worker.test.exe
Build ID: C:\Users\Yamabaka\AppData\Local\Temp\go-build1478970430\b001\worker.test.exe2025-11-01 22:46:45.7195369 +0300 MSK
Type: alloc_space
Time: 2025-11-01 17:32:47 MSK
Showing nodes accounting for 17.53MB, 1.33% of 1318.58MB total
Dropped 3 nodes (cum <= 6.59MB)
      flat  flat%   sum%        cum   cum%
   16.52MB  1.25%  1.25%    16.52MB  1.25%  github.com/Okenamay/shorturl.git/internal/worker.softDeleter.func1
   -3.50MB  0.27%  0.99%    -3.50MB  0.27%  github.com/Okenamay/shorturl.git/internal/worker.BenchmarkSendToDelete.func2
       3MB  0.23%  1.22%    19.02MB  1.44%  github.com/Okenamay/shorturl.git/internal/worker.softDeleter.func2
       2MB  0.15%  1.37%        2MB  0.15%  runtime.allocm
    0.50MB 0.038%  1.41%     0.50MB 0.038%  unicode.map.init.1
    0.50MB 0.038%  1.44%     0.50MB 0.038%  golang.org/x/text/internal/language.map.init.1
   -0.50MB 0.038%  1.41%    -0.50MB 0.038%  regexp/syntax.(*compiler).inst (inline)
   -0.50MB 0.038%  1.37%    -0.50MB 0.038%  compress/flate.newHuffmanEncoder (inline)
   -0.50MB 0.038%  1.33%    -0.50MB 0.038%  time.NewTicker
         0     0%  1.33%    -0.50MB 0.038%  compress/flate.generateFixedLiteralEncoding
         0     0%  1.33%    -0.50MB 0.038%  compress/flate.init
         0     0%  1.33%     0.50MB 0.038%  golang.org/x/text/internal/language.init
         0     0%  1.33%    -0.50MB 0.038%  gopkg.in/yaml%2ev3.init
         0     0%  1.33%    -0.50MB 0.038%  regexp.Compile (inline)
         0     0%  1.33%    -0.50MB 0.038%  regexp.MustCompile
         0     0%  1.33%    -0.50MB 0.038%  regexp.compile
         0     0%  1.33%    -0.50MB 0.038%  regexp/syntax.(*compiler).cap (inline)
         0     0%  1.33%    -0.50MB 0.038%  regexp/syntax.(*compiler).compile
         0     0%  1.33%    -0.50MB 0.038%  regexp/syntax.Compile
         0     0%  1.33%        2MB  0.15%  runtime.mcall
         0     0%  1.33%        2MB  0.15%  runtime.newm
         0     0%  1.33%        2MB  0.15%  runtime.park_m
         0     0%  1.33%        2MB  0.15%  runtime.resetspinning
         0     0%  1.33%        2MB  0.15%  runtime.schedule
         0     0%  1.33%        2MB  0.15%  runtime.startm
         0     0%  1.33%        2MB  0.15%  runtime.wakep
         0     0%  1.33%    -3.50MB  0.27%  testing.(*B).RunParallel.func1
         0     0%  1.33%     0.50MB 0.038%  unicode.init

==================================================
Comparison finished.