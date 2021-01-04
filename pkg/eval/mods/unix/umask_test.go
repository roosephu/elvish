// +build !windows,!plan9,!js

package unix

import (
	"testing"

	"github.com/elves/elvish/pkg/eval"
	"github.com/elves/elvish/pkg/eval/errs"
	. "github.com/elves/elvish/pkg/eval/evaltest"
)

// Note that this unit test assumes a UNIX environment with a POSIX compatible
// /bin/sh program.
func TestUmask(t *testing.T) {
	setup := func(ev *eval.Evaler) {
		ev.SetGlobal(eval.NsBuilder{}.AddNs("unix", Ns).Ns())
	}
	TestWithSetup(t, setup,
		// We have to start with a known umask value.
		That(`unix:umask = 022`).Puts(),
		That(`put $unix:umask`).Puts(`0o022`),
		// Verify that mutating the value and outputting the new value works.
		That(`unix:umask = 23`).Puts(),
		That(`put $unix:umask`).Puts(`0o023`),
		That(`unix:umask = 0o75`).Puts(),
		That(`put $unix:umask`).Puts(`0o075`),
		// Verify that a temporary umask change is reverted upon completion of
		// the command. Both for builtin and external commands.
		That(`unix:umask=012 put $unix:umask`).Puts(`0o012`),
		That(`unix:umask=0o23 /bin/sh -c 'umask'`).Prints("0023\n"),
		That(`unix:umask=56 /bin/sh -c 'umask'`).Prints("0056\n"),
		That(`put $unix:umask`).Puts(`0o075`),
		// People won't normally use non-octal bases but make sure these cases
		// behave sensibly given that Elvish supports number literals with an
		// explicit base.
		That(`unix:umask=0x43 /bin/sh -c 'umask'`).Prints("0103\n"),
		That(`unix:umask=0b001010100 sh -c 'umask'`).Prints("0124\n"),
		// We should be back to our expected umask given the preceding tests
		// applied a temporary change to that process attribute.
		That(`put $unix:umask`).Puts(`0o075`),
		// An explicit float64 value should be handled correctly.
		That(`unix:umask=(float64 0o17) put $unix:umask`).Puts(`0o017`),
		That(`unix:umask=(float64 123.4)`).Throws(errs.BadValue{
			What: "umask", Valid: validUmaskMsg, Actual: "123.4"}),

		// An invalid string should raise the expected exception.
		That(`unix:umask=022z`).Throws(errs.BadValue{
			What: "umask", Valid: validUmaskMsg, Actual: "022z"}),

		// An invalid data type should raise the expected exception.
		That(`unix:umask=[1]`).Throws(errs.BadValue{
			What: "umask", Valid: validUmaskMsg, Actual: "[1]"}),

		// Values outside the legal range should raise the expected exception.
		//
		// TODO: Switch to `%O` when Go 1.15 is the minimum acceptable version.
		// Until then the formatting of negative numbers will be weird.
		That(`unix:umask=0o1000`).Throws(errs.OutOfRange{
			What: "umask", ValidLow: "0", ValidHigh: "0o777", Actual: "0o1000"}),

		That(`unix:umask=-1`).Throws(errs.OutOfRange{
			What: "umask", ValidLow: "0", ValidHigh: "0o777", Actual: "0o-1"}),
	)
}
