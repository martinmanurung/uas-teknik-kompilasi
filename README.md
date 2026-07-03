# Dokumentasi Proyek UAS Teknik Kompilasi

### Identitas

- **Nama** : MARTIN RIVALDO MANURUNG
- **NIM** : 231011400295
- **Mata Kuliah** : Teknik Kompilasi
- **Dosen Pengampu** : AGUNG PERDANANTO S.Kom, M.Kom

## 1. Gambaran Umum

Proyek ini mensimulasikan tahapan utama proses kompilasi untuk konstruksi **deklarasi fungsi/metode**. Implementasi dibuat menggunakan bahasa **Go** dan menampilkan empat tahap utama:

1. Analisis leksikal
2. Analisis sintaksis
3. Analisis semantik
4. Generasi kode antara atau Three-Address Code (TAC)

Fokus utama program adalah membaca sebuah source code sederhana, memecahnya menjadi token, membentuk AST sederhana, melakukan validasi tipe dasar, lalu menghasilkan TAC yang mewakili isi fungsi.

## 2. Konstruksi yang Dipilih

Konstruksi yang diimplementasikan adalah **deklarasi fungsi/metode**.

Contoh bentuk input yang didukung:

```text
func hitungTotal ( a int , b int ) int { c = a + b return c }
```

Struktur tersebut merepresentasikan:

- nama fungsi: `hitungTotal`
- parameter: `a int`, `b int`
- return type: `int`
- body: `c = a + b return c`

## 3. Pattern / Grammar

Pola sintaks yang digunakan pada program ini dapat ditulis dalam bentuk BNF sederhana sebagai berikut:

```text
<func_decl>  ::= "func" <identifier> "(" <params> ")" <type> "{" <body> "}"
<params>     ::= <param> | <param> "," <params>
<param>      ::= <identifier> <type>
<type>       ::= "int" | "float" | "string" | "void"
<body>       ::= <statement> | <statement> <body>
<statement>  ::= <assignment> | <return_stmt>
<assignment> ::= <identifier> "=" <expression>
<return_stmt>::= "return" <identifier>
<expression> ::= <identifier> <operator> <identifier> | <identifier>
<operator>   ::= "+" | "-" | "*" | "/"
```

Grammar ini dibuat sederhana agar sesuai dengan tujuan tugas, yaitu menunjukkan alur kompilasi secara konseptual, bukan membangun compiler penuh.

## 4. Struktur Program

File utama proyek adalah [main.go](main.go). Di dalamnya terdapat beberapa bagian penting:

- `Token`: menyimpan hasil analisis leksikal
- `ASTNode`: menyimpan simpul AST
- `FuncDeclCompiler`: struct utama yang menyimpan source code, daftar token, posisi parsing, dan counter temporary variable
- `LexicalAnalysis()`: melakukan tokenisasi
- `SyntaxAnalysis()`: membentuk AST
- `SemanticAnalysis()`: memeriksa validitas tipe data
- `GenerateTAC()`: menghasilkan TAC
- `main()`: menjalankan seluruh tahapan dan menampilkan hasil

## 5. Penjelasan Tiap Tahap

### 5.1 Analisis Leksikal

Tahap leksikal ada di fungsi `LexicalAnalysis()`.

Cara kerjanya:

1. Source code dibersihkan dengan menambahkan spasi di sekitar simbol penting seperti `(`, `)`, `{`, `}`, `,`, `=`, `+`, `-`, `*`, dan `/`.
2. String kemudian dipecah menggunakan `strings.Fields()`.
3. Setiap hasil pecahan dimasukkan ke dalam `Token`.

Hasil tahap ini adalah daftar token yang nantinya dipakai untuk parsing.

### 5.2 Analisis Sintaksis

Tahap sintaksis ada di fungsi `SyntaxAnalysis()`.

Program memeriksa apakah token pertama adalah `func`. Jika iya, parser membaca:

- nama fungsi
- daftar parameter di dalam tanda kurung
- return type setelah `)`
- isi body di dalam `{ }`

Setelah itu parser membentuk AST sederhana dengan akar `FuncDecl` dan tiga anak node:

- `Params`
- `ReturnType`
- `Body`

AST ini belum lengkap seperti AST compiler nyata, tetapi sudah cukup untuk menunjukkan struktur hierarki sintaks.

### 5.3 Analisis Semantik

Tahap semantik ada di fungsi `SemanticAnalysis()`.

Pemeriksaan yang dilakukan adalah:

- memastikan node AST yang diterima bertipe `FuncDecl`
- memvalidasi return type terhadap daftar tipe yang diizinkan: `int`, `float`, `string`, dan `void`
- memvalidasi tipe parameter dengan daftar tipe yang sama

Jika semua valid, program menampilkan pesan bahwa semantic analysis berhasil.

### 5.4 Generasi TAC

Tahap generasi kode antara ada di fungsi `GenerateTAC()`.

Program menghasilkan TAC dengan format sederhana:

- `BeginFunc <nama_fungsi>`
- `PopParam <nama_parameter>` untuk setiap parameter
- instruksi untuk assignment atau operasi aritmetika
- `Return <nilai>`
- `EndFunc <nama_fungsi>`

Untuk ekspresi seperti `a + b`, program membuat temporary variable terlebih dahulu, misalnya:

```text
t1 = a + b
c = t1
```

## 6. Contoh Output Program

Input yang digunakan pada `main()`:

```text
func hitungTotal ( a int , b int ) int { c = a + b return c }
```

Output yang dihasilkan program:

```text
--- 1. Analisis Leksikal (Tokens) ---
'func' 'hitungTotal' '(' 'a' 'int' ',' 'b' 'int' ')' 'int' '{' 'c' '=' 'a' '+' 'b' 'return' 'c' '}'

--- 2. Analisis Sintaksis (Abstract Syntax Tree) ---
FuncDecl: hitungTotal
	├─ Params: a int, b int
	├─ ReturnType: int
	└─ Body: c = a + b return c

--- 3. Analisis Semantik ---
Semantic Analysis Passed: Tipe data parameter dan return type valid.

--- 4. Generasi Kode Antara (TAC) ---
BeginFunc hitungTotal
PopParam a
PopParam b
t1 = a + b
c = t1
Return c
EndFunc hitungTotal
```

## 7. Cara Menjalankan Program

Pastikan sudah berada di folder proyek `uas-teknik-kompilasi`, lalu jalankan:

```bash
go run main.go
```

Jika ingin membuat binary terlebih dahulu:

```bash
go build -o uas-kompilasi
./uas-kompilasi
```

## 8. Kesesuaian dengan Ketentuan UAS

Dokumen tugas UAS meminta adanya:

- pemilihan satu konstruksi sintaksis
- pattern / grammar
- implementasi leksikal, sintaksis, semantik, dan TAC
- penjelasan lengkap dalam format Markdown atau PDF

Proyek ini sudah memenuhi keempat poin tersebut dengan konstruksi **deklarasi fungsi/metode**.

## 9. Catatan dan Keterbatasan

Implementasi ini bersifat simulasi dan sederhana. Beberapa keterbatasannya adalah:

- parsing hanya menangani format deklarasi fungsi yang sangat spesifik
- body fungsi dibaca sebagai deretan token linear, belum menjadi statement tree yang lengkap
- validasi semantik masih terbatas pada pengecekan tipe data dasar
- TAC yang dihasilkan hanya mendukung assignment sederhana, operasi aritmetika dasar, dan return

Walaupun sederhana, struktur ini sudah cukup untuk menunjukkan alur dasar proses kompilasi sesuai kebutuhan tugas.
