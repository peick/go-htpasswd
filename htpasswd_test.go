package htpasswd

import (
	_ "embed"
	"os"
	"strings"
	"testing"
)

type testUser struct {
	username string
	password string
}

var testUsers = []testUser{
	{"user1", "mickey5"},
	{"user2", "alexandrew"},
	{"user3", "hawaiicats78"},
	{"user4", "DIENOW"},
	{"user5", "e8f685"},
	{"user6", "Rickygirl03"},
	{"user7", "123vb123"},
	{"user8", "sheng060576"},
	{"user9", "hansisme"},
	{"user10", "h4ck3rs311t3"},
	{"user11", "K90JyTGA"},
	{"user12", "aspire5101"},
	{"user13", "553568"},
	{"user14", "SRI"},
	{"user15", "maxmus"},
	{"user16", "a5xp9707"},
	{"user17", "tomasrim"},
	{"user18", "2a0mag"},
	{"user19", "wmsfht"},
	{"user20", "webmaster2364288"},
	{"user21", "121516m"},
	{"user22", "T69228803"},
	{"user23", "qq820221"},
	{"user24", "chenfy"},
	{"user25", "www.debure.net"},
	{"user26", "1333e763"},
	{"user27", "burberries"},
	{"user28", "chanmee14"},
	{"user29", "65432106543210"},
	{"user30", "powernet"},
	{"user31", "a2d8i6a7"},
	{"user32", "gvs9ptc"},
	{"user33", "Pookie"},
	{"user34", "lorissss"},
	{"user35", "ess"},
	{"user36", "sparra"},
	{"user37", "allysson"},
	{"user38", "99128008"},
	{"user39", "evisanne"},
	{"user40", "qfxg7x9l"},
	{"user41", "03415"},
	{"user42", "87832309"},
	{"user43", "816283"},
	{"user44", "banach12"},
	{"user45", "sjdszpsc"},
	{"user46", "changsing"},
	{"user47", "56339388"},
	{"user48", "52114157"},
	{"user49", "jinebimb"},
	{"user50", "erol43"},
	{"user51", "2yagos"},
	{"user52", "habparty!"},
	{"user53", "tangjianhui"},
	{"user54", "serandah"},
	{"user55", "mirrages"},
	{"user56", "mantgaxxl"},
	{"user57", "45738901"},
	{"user58", "g523minna"},
	{"user59", "j202020"},
	{"user60", "g@mmaecho"},
	{"user61", "042380"},
	{"user62", "ASRuin"},
	{"user63", "061990"},
	{"user64", "ysoline"},
	{"user65", "liuzhouzhou"},
	{"user66", "b0000000wind"},
	{"user67", "7913456852"},
	{"user68", "9008"},
	{"user69", "waitlin11"},
	{"user70", "8fdakar"},
	{"user71", "eisball"},
	{"user72", "jenna17"},
	{"user73", "belkadonam"},
	{"user74", "tfyuj9JW"},
	{"user75", "nihaijidema"},
	{"user76", "talapia"},
	{"user77", "7376220"},
	{"user78", "c7m8e1xsc3"},
	{"user79", "84129793"},
	{"user80", "test1000"},
	{"user81", "ecmanhatten"},
	{"user82", "EvanYo3327"},
	{"user83", "269john139"},
	{"user84", "3348159zw"},
	{"user85", "lu184020"},
	{"user86", "aszasw"},
	{"user87", "33059049"},
	{"user88", "li3255265"},
	{"user89", "kerrihayes"},
	{"user90", "0167681809"},
	{"user91", "stefano123"},
	{"user92", "15054652730"},
	{"user93", "natdvd213"},
	{"user94", "680929"},
	{"user95", "steelpad8"},
	{"user96", "374710"},
	{"user97", "394114"},
	{"user98", "24347"},
	{"user99", "krait93"},
	{"user100", "5164794"},
	{"user101", "rswCyJE5"},
	{"user102", "31480019"},
	{"user103", "19830907ok"},
	{"user104", "zlsmhzlsmh"},
	{"user105", "Zengatsu"},
	{"user106", "0127603331"},
	{"user107", "axelle77"},
	{"user108", "password2147"},
	{"user109", "olixkl8b"},
	{"user110", "maiwen"},
	{"user111", "198613"},
	{"user112", "s17kr8wu"},
	{"user113", "biker02"},
	{"user114", "m1399"},
	{"user115", "a2dc6a"},
	{"user116", "zhd8902960"},
	{"user117", "parasuta"},
	{"user118", "the1secret"},
	{"user119", "teddy14"},
	{"user120", "4516388amt"},
	{"user121", "245520"},
	{"user122", "D34dw00d"},
	{"user123", "officiel"},
	{"user124", "36653665"},
	{"user125", "hipol"},
	{"user126", "Nylon0"},
	{"user127", "caitlyne6"},
	{"user128", "dogzilla"},
	{"user129", "lemegaboss"},
	{"user130", "c0valerius"},
	{"user131", "liseczek44"},
	{"user132", "saulosi"},
	{"user133", "53522"},
	{"user134", "ajgebam"},
	{"user135", "freshplayer"},
	{"user136", "logistica1"},
	{"user137", "12calo66"},
	{"user138", "kenno"},
	{"user139", "34639399"},
	{"user140", "0408636405"},
	{"user141", "weezer12"},
	{"user142", "9888735777"},
	{"user143", "7771877"},
	{"user144", "6620852"},
	{"user145", "98billiards"},
	{"user146", "angelik"},
	{"user147", "86815057"},
	{"user148", "p16alfalfa"},
	{"user149", "7236118"},
	{"user150", "glock17l"},
	{"user151", "sigmundm"},
	{"user152", "ltbgeqsd"},
	{"user153", "wqnd8k2m"},
	{"user154", "yangjunjie"},
	{"user155", "manjinder"},
	{"user156", "nick2000"},
	{"user157", "193416"},
	{"user158", "pang168"},
	{"user159", "454016"},
	{"user160", "phair08"},
	{"user161", "10252007cw"},
	{"user162", "zhuzhuzhu"},
	{"user163", "metafunds"},
	{"user164", "smash"},
	{"user165", "76387638"},
	{"user166", "S226811954"},
	{"user167", "mintymoo00"},
	{"user168", "seven711"},
	{"user169", "924414"},
	{"user170", "changchengxu"},
	{"user171", "alaska58"},
	{"user172", "7678208"},
	{"user173", "szazsoo73"},
	{"user174", "3830371"},
	{"user175", "0qdzx66b"},
	{"user176", "09124248099"},
	{"user177", "bachrain"},
	{"user178", "sJsSdFBY"},
	{"user179", "676215000"},
	{"user180", "nimamapwoaini"},
	{"user181", "nitsuj"},
	{"user182", "cukierek2003"},
	{"user183", "seeder"},
	{"user184", "00167148786"},
	{"user185", "ashok198"},
	{"user186", "kt2116"},
	{"user187", "another82"},
	{"user188", "75995794"},
	{"user189", "19901130"},
	{"user190", "gijs010389"},
	{"user191", "26263199"},
	{"user192", "hi1j42x8"},
	{"user193", "6922235"},
	{"user194", "67749330"},
	{"user195", "ccpatrik"},
	{"user196", "summer3011"},
	{"user197", "331516"},
	{"user198", "135745"},
	{"user199", "603762004"},
	{"user200", "29011985"},
	{"user201", "29011985"},
}

//go:embed testdata/htpasswd/textPlain
var textPlain string

//go:embed testdata/htpasswd/textApr1
var textApr1 string

//go:embed testdata/htpasswd/textSha
var textSha string

//go:embed testdata/htpasswd/textBcrypt
var textBcrypt string

//go:embed testdata/htpasswd/textSsha
var textSsha string

//go:embed testdata/htpasswd/textMd5Crypt
var textMd5Crypt string

//go:embed testdata/htpasswd/testCryptSha256
var testCryptSha256 string

//go:embed testdata/htpasswd/testCryptSha512
var testCryptSha512 string

func testSystemReader(t *testing.T, name string, contents string) {
	r := strings.NewReader(contents)

	htp, err := NewFromReader(r)
	if err != nil {
		t.Fatalf("Failed to read htpasswd reader")
	}

	for _, u := range testUsers {
		if good := htp.Match(u.username, u.password); !good {
			t.Errorf("%s user %s, password %s failed to authenticate: %t", name, u.username, u.password, good)
		}

		notPass := u.password + "not"
		if bad := htp.Match(u.username, notPass); bad {
			t.Errorf("%s user %s, password %s erroneously authenticated: %t", name, u.username, notPass, bad)
		}
	}
}

func testSystem(t *testing.T, name string, contents string) {
	f, err := os.CreateTemp("", "gohtpasswd")
	if err != nil {
		t.Fatalf("Failed to make temp file: %s", err.Error())
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(contents); err != nil {
		t.Fatalf("Failed to write temporary file: %s", err.Error())
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %s", err.Error())
	}

	htp, err := New(f.Name())
	if err != nil {
		t.Fatalf("Failed to read htpasswd file")
	}

	for _, u := range testUsers {

		if good := htp.Match(u.username, u.password); !good {
			t.Errorf("%s user %s, password %s failed to authenticate: %t", name, u.username, u.password, good)
		}

		notPass := u.password + "not"
		if bad := htp.Match(u.username, notPass); bad {
			t.Errorf("%s user %s, password %s erroneously authenticated: %t", name, u.username, notPass, bad)
		}
	}
}

func Test_PlainReader(t *testing.T) { testSystemReader(t, "plain", textPlain) }
func Test_PlainFile(t *testing.T)   { testSystem(t, "plain", textPlain) }

func Test_ShaReader(t *testing.T) { testSystemReader(t, "sha", textSha) }
func Test_ShaFile(t *testing.T)   { testSystem(t, "sha", textSha) }

func Test_Apr1Reader(t *testing.T) { testSystemReader(t, "md5", textApr1) }
func Test_Apr1File(t *testing.T)   { testSystem(t, "md5", textApr1) }

func Test_Md5Reader(t *testing.T) { testSystemReader(t, "md5", textMd5Crypt) }
func Test_Md5File(t *testing.T)   { testSystem(t, "md5", textMd5Crypt) }

func Test_BcryptReader(t *testing.T) { testSystemReader(t, "bcrypt", textBcrypt) }
func Test_BcryptFile(t *testing.T)   { testSystem(t, "bcrypt", textBcrypt) }

func Test_SshaReader(t *testing.T) { testSystemReader(t, "ssha", textSsha) }
func Test_SshaFile(t *testing.T)   { testSystem(t, "ssha", textSsha) }

func Test_CryptSha256Reader(t *testing.T) { testSystemReader(t, "crypt-sha256", testCryptSha256) }
func Test_CryptSha256File(t *testing.T)   { testSystem(t, "crypt-sha256", testCryptSha256) }

func Test_CryptSha512Reader(t *testing.T) { testSystemReader(t, "crypt-sha512", testCryptSha512) }
func Test_CryptSha512File(t *testing.T)   { testSystem(t, "crypt-sha512", testCryptSha512) }
