package atom

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEmbeddedAtom(t *testing.T) {
	Convey("create embedded atom", t, func() {
		Convey("atom `test`", func() {
			a := embedAtom([]byte("test"))

			So(uint32(a), ShouldEqual, 0xf4657374)
			So(a.String(), ShouldEqual, "test")
			So(a.Len(), ShouldEqual, 4)
			So(a.IsEmbedded(), ShouldBeTrue)
		})

		Convey("atom `abc`", func() {
			a := embedAtom([]byte("abc"))

			So(uint32(a), ShouldEqual, 0x83616263)
			So(a.String(), ShouldEqual, "abc")
			So(a.Len(), ShouldEqual, 3)
			So(a.IsEmbedded(), ShouldBeTrue)
		})

		Convey("atom `ok`", func() {
			a := embedAtom([]byte("ok"))

			So(uint32(a), ShouldEqual, 0x82006f6b)
			So(a.String(), ShouldEqual, "ok")
			So(a.Len(), ShouldEqual, 2)
			So(a.IsEmbedded(), ShouldBeTrue)
		})

		Convey("atom `I`", func() {
			a := embedAtom([]byte("I"))

			So(uint32(a), ShouldEqual, 0x81000049)
			So(a.String(), ShouldEqual, "I")
			So(a.Len(), ShouldEqual, 1)
			So(a.IsEmbedded(), ShouldBeTrue)
		})

		Convey("atom `我`", func() {
			s := "我"
			a := embedAtom([]byte(s))

			So(uint32(a), ShouldEqual, 0x83e68891)
			So(a.String(), ShouldEqual, s)
			So(a.Len(), ShouldEqual, 3)
			So(a.IsEmbedded(), ShouldBeTrue)

			So(Lookup(s), ShouldEqual, a)
		})

		Convey("atom `too long`", func() {
			So(embedAtom([]byte("too long")), ShouldEqual, Empty)
		})
	})
}

func TestAddAtom(t *testing.T) {
	Convey("add some atom", t, func() {
		Load(nil, nil)

		Convey("add `golang` atom", func() {
			s := []byte("golang")

			So(findAtomInCache(s), ShouldEqual, Empty)
			So(findAtomInData(s), ShouldEqual, Empty)

			a := addAtom(s)

			So(uint32(a), ShouldEqual, 0x6000000)
			So(a.String(), ShouldEqual, "golang")
			So(a.Len(), ShouldEqual, 6)
			So(a.IsEmbedded(), ShouldBeFalse)

			So(findAtomInCache(s), ShouldEqual, a)
			So(findAtomInData(s), ShouldEqual, a)
		})

		Convey("add `汉字` atom", func() {
			s := "汉字"
			a := New(s)

			So(uint32(a), ShouldEqual, 0x6000000)
			So(a.String(), ShouldResemble, s)
			So(a.Len(), ShouldEqual, 6)
			So(a.IsEmbedded(), ShouldBeFalse)

			So(findAtomInCache([]byte(s)), ShouldEqual, a)
			So(findAtomInData([]byte(s)), ShouldEqual, a)

			Convey("lookup it should return the atom", func() {
				So(Lookup(s), ShouldEqual, a)
			})

			Convey("add it again should return same atom", func() {
				So(New(s), ShouldEqual, a)
			})
		})

		Convey("add a random atom", func() {
			s := make([]byte, MaxAtomLen)

			for i := 0; i < MaxAtomLen; i++ {
				s[i] = byte(rand.Uint32())
			}

			a := addAtom(s)

			So(uint32(a), ShouldEqual, 0x7f000000)
			So(a.Bytes(), ShouldResemble, s)
			So(a.Len(), ShouldEqual, len(s))
			So(a.IsEmbedded(), ShouldBeFalse)

			So(findAtomInCache(s), ShouldEqual, a)
			So(findAtomInData(s), ShouldEqual, a)
		})
	})
}

func TestInvalidAtom(t *testing.T) {
	Convey("add some invalid atom", t, func() {
		Convey("add a empty atom", func() {
			So(newAtom(0, 0), ShouldEqual, Empty)

			a := New("")

			So(uint32(a), ShouldEqual, 0)
			So(a.IsEmpty(), ShouldBeTrue)
			So(a, ShouldEqual, Empty)
			So(a.Bytes(), ShouldBeNil)
			So(a.String(), ShouldEqual, "")
		})

		Convey("add a too long atom", func() {
			So(newAtom(0, MaxAtomLen+1), ShouldEqual, Empty)

			So(New(string(make([]byte, 128))), ShouldEqual, Empty)
		})

		Convey("add a atom out of range", func() {
			So(newAtom(atomData.Len(), 0), ShouldEqual, Empty)
		})

		Convey("lookup a empty atom", func() {
			So(Lookup(""), ShouldEqual, Empty)
		})
	})
}
