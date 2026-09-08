//go:build amd64 && !purego

#include "textflag.h"

// Four pixels at a time in SSE2, matching Blend exactly.
//
// The arithmetic is the scalar one rearranged so that two channels share a
// register. A pixel splits into two pairs of 16-bit fields — red with blue,
// and alpha with green — and each pair is multiplied, added and divided as
// one. That is what makes four pixels cost roughly what one costs scalar.
//
// The division by 255 is exact rather than approximate: for any t below
// 65536, (t + (t>>8) + 1) >> 8 is t/255, and t here cannot exceed 65152
// because the two weights always sum to 255.

// rbMask isolates red and blue, and doubles as 255 in every 16-bit field —
// the same bit pattern serves both.
DATA rbMask<>+0(SB)/8, $0x00FF00FF00FF00FF
DATA rbMask<>+8(SB)/8, $0x00FF00FF00FF00FF
GLOBL rbMask<>(SB), RODATA|NOPTR, $16

DATA round<>+0(SB)/8, $0x007F007F007F007F
DATA round<>+8(SB)/8, $0x007F007F007F007F
GLOBL round<>(SB), RODATA|NOPTR, $16

DATA one<>+0(SB)/8, $0x0001000100010001
DATA one<>+8(SB)/8, $0x0001000100010001
GLOBL one<>(SB), RODATA|NOPTR, $16

DATA opaque<>+0(SB)/8, $0xFF000000FF000000
DATA opaque<>+8(SB)/8, $0xFF000000FF000000
GLOBL opaque<>(SB), RODATA|NOPTR, $16

// func Rows(dst, src *uint32, count int)
TEXT ·Rows(SB), NOSPLIT, $0-24
	MOVQ dst+0(FP), DI
	MOVQ src+8(FP), SI
	MOVQ count+16(FP), CX
	SHRQ $2, CX               // groups of four
	JZ   done

	MOVOU rbMask<>(SB), X11
	MOVOU round<>(SB), X12
	MOVOU one<>(SB), X13
	MOVOU opaque<>(SB), X14
	PXOR  X10, X10

loop:
	MOVOU (SI), X0            // s
	MOVOU (DI), X1            // d

	MOVOU X0, X2
	PSRLL $24, X2             // a, in the low byte of each 32-bit lane

	MOVOU X2, X3
	PSLLL $16, X3
	POR   X2, X3              // av: alpha in both 16-bit fields

	MOVOU X11, X4
	PSUBW X3, X4              // iv = 255 - av

	MOVOU X0, X5
	PAND  X11, X5             // src red and blue
	MOVOU X0, X6
	PSRLL $8, X6
	PAND  X11, X6             // src alpha and green
	MOVOU X1, X7
	PAND  X11, X7             // dst red and blue
	MOVOU X1, X8
	PSRLL $8, X8
	PAND  X11, X8             // dst alpha and green

	PMULLW X3, X5             // srb * a
	PMULLW X4, X7             // drb * (255-a)
	PADDW  X7, X5
	PADDW  X12, X5            // + 127
	MOVOU  X5, X9
	PSRLW  $8, X9
	PADDW  X9, X5
	PADDW  X13, X5
	PSRLW  $8, X5             // red and blue, divided by 255

	PMULLW X3, X6
	PMULLW X4, X8
	PADDW  X8, X6
	PADDW  X12, X6
	MOVOU  X6, X9
	PSRLW  $8, X9
	PADDW  X9, X6
	PADDW  X13, X6
	PSRLW  $8, X6             // alpha and green, divided by 255

	PSLLL $8, X6              // green and alpha back into place
	POR   X6, X5
	POR   X14, X5             // the result is opaque, as Blend's is

	// Where the source alpha was zero, Blend gives the destination back
	// untouched — alpha included — so those lanes are selected out again.
	PCMPEQL X10, X2           // all ones where a == 0
	MOVOU   X2, X9
	PAND    X1, X9
	PANDN   X5, X2
	POR     X9, X2
	MOVOU   X2, (DI)

	ADDQ $16, SI
	ADDQ $16, DI
	DECQ CX
	JNZ  loop

done:
	RET

// func Solid(dst *uint32, colour uint32, count int)
TEXT ·Solid(SB), NOSPLIT, $0-24
	MOVQ dst+0(FP), DI
	MOVL colour+8(FP), AX
	MOVQ count+16(FP), CX
	SHRQ $2, CX
	JZ   solidDone

	MOVOU rbMask<>(SB), X11
	MOVOU round<>(SB), X12
	MOVOU one<>(SB), X13
	MOVOU opaque<>(SB), X14

	// One colour for every pixel, so everything that depends only on the
	// source is worked out once and kept in a register.
	MOVD  AX, X0
	PSHUFL $0x00, X0, X0      // broadcast to all four lanes

	MOVOU X0, X2
	PSRLL $24, X2
	MOVOU X2, X3
	PSLLL $16, X3
	POR   X2, X3              // av

	MOVOU X11, X4
	PSUBW X3, X4              // iv

	MOVOU X0, X5
	PAND  X11, X5
	PMULLW X3, X5             // srb * a, constant
	MOVOU X0, X6
	PSRLL $8, X6
	PAND  X11, X6
	PMULLW X3, X6             // sag * a, constant

solidLoop:
	MOVOU (DI), X1            // d

	MOVOU X1, X7
	PAND  X11, X7
	MOVOU X1, X8
	PSRLL $8, X8
	PAND  X11, X8

	PMULLW X4, X7
	PADDW  X5, X7
	PADDW  X12, X7
	MOVOU  X7, X9
	PSRLW  $8, X9
	PADDW  X9, X7
	PADDW  X13, X7
	PSRLW  $8, X7

	PMULLW X4, X8
	PADDW  X6, X8
	PADDW  X12, X8
	MOVOU  X8, X9
	PSRLW  $8, X9
	PADDW  X9, X8
	PADDW  X13, X8
	PSRLW  $8, X8

	PSLLL $8, X8
	POR   X8, X7
	POR   X14, X7
	MOVOU X7, (DI)

	ADDQ $16, DI
	DECQ CX
	JNZ  solidLoop

solidDone:
	RET
