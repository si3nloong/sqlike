package charset

// Code :
type Code string

// codes :
const (
	Utf8mb4  Code = "utf8mb4"  // UTF-8 Unicode
	Utf8     Code = "utf8"     // UTF-8 Unicode
	Utf16    Code = "utf16"    // UTF-16 Unicode
	Utf32    Code = "utf32"    // UTF-32 Unicode
	Latin1   Code = "latin1"   // cp1252 West European
	Latin2   Code = "latin2"   // ISO 8859-2 Central European
	Latin5   Code = "latin5"   // ISO 8859-9 Turkish
	Latin7   Code = "latin7"   // ISO 8859-13 Baltic
	Big5     Code = "big5"     // Big5 Traditional Chinese
	ASCII    Code = "ascii"    // US ASCII
	ARMSCII8 Code = "armscii8" // ARMSCII-8 Armenian
	Greek    Code = "greek"    // ISO 8859-7 Greek
	HP8      Code = "hp8"      // HP West European
	KOI8R    Code = "koi8r"    // KOI8-R Relcom Russian
	SWE7     Code = "swe7"     // 7bit Swedish
	UJIS     Code = "ujis"     // EUC-JP Japanese
	SJIS     Code = "sjis"     // Shift-JIS Japanese
	TIS620   Code = "tis620"   // TIS620 Thai
	EUCKR    Code = "euckr"    // EUC-KR Korean
	GBK      Code = "gbk"      // GBK Simplified Chinese
	KEYBCS2  Code = "keybcs2"  // DOS Kamenicky Czech-Slovak
	GB2312   Code = "gb2312"   // GB2312 Simplified Chinese
	CP866    Code = "cp866"    // DOS Russian
	CP932    Code = "cp932"    // SJIS for Windows Japanese
	CP1250   Code = "cp1250"   // Windows Central European
	CP1251   Code = "cp1251"   // Windows Cyrillic
	CP1256   Code = "cp1256"   // Windows Arabic
	CP1257   Code = "cp1257"   // Windows Baltic
	Binary   Code = "binary"   // Binary pseudo charset
	EUCJPMS  Code = "eucjpms"  // UJIS for Windows Japanese
)
