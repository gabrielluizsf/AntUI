package sdk

// Ship is what goes to the store: 64-bit ARM is required, and 32-bit ARM is
// what keeps the app installable on a phone from before 2016. The x86 pair
// is for the emulator and is deliberately not here — including it in a
// bundle costs every user the download of a binary no phone can run.
var Ship = []ABI{Arm64, Arm}
