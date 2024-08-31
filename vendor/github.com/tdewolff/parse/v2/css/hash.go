package css

// generated by hasher -type=Hash -file=hash.go; DO NOT EDIT, except for adding more constants to the list and rerun go generate

// uses github.com/tdewolff/hasher
//go:generate hasher -type=Hash -file=hash.go

// Hash defines perfect hashes for a predefined list of strings
type Hash uint32

// Unique hash definitions to be used instead of strings
const (
	Document  Hash = 0x8    // document
	Font_Face Hash = 0x809  // font-face
	Keyframes Hash = 0x1109 // keyframes
	Media     Hash = 0x2105 // media
	Page      Hash = 0x2604 // page
	Supports  Hash = 0x1908 // supports
)

// String returns the hash' name.
func (i Hash) String() string {
	start := uint32(i >> 8)
	n := uint32(i & 0xff)
	if start+n > uint32(len(_Hash_text)) {
		return ""
	}
	return _Hash_text[start : start+n]
}

// ToHash returns the hash whose name is s. It returns zero if there is no
// such hash. It is case sensitive.
func ToHash(s []byte) Hash {
	if len(s) == 0 || len(s) > _Hash_maxLen {
		return 0
	}
	h := uint32(_Hash_hash0)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	if i := _Hash_table[h&uint32(len(_Hash_table)-1)]; int(i&0xff) == len(s) {
		t := _Hash_text[i>>8 : i>>8+i&0xff]
		for i := 0; i < len(s); i++ {
			if t[i] != s[i] {
				goto NEXT
			}
		}
		return i
	}
NEXT:
	if i := _Hash_table[(h>>16)&uint32(len(_Hash_table)-1)]; int(i&0xff) == len(s) {
		t := _Hash_text[i>>8 : i>>8+i&0xff]
		for i := 0; i < len(s); i++ {
			if t[i] != s[i] {
				return 0
			}
		}
		return i
	}
	return 0
}

const _Hash_hash0 = 0x9acb0442
const _Hash_maxLen = 9
const _Hash_text = "documentfont-facekeyframesupportsmediapage"

var _Hash_table = [1 << 3]Hash{
	0x1: 0x2604, // page
	0x2: 0x2105, // media
	0x3: 0x809,  // font-face
	0x5: 0x1109, // keyframes
	0x6: 0x1908, // supports
	0x7: 0x8,    // document
}
