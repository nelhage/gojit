// textflag.h
// Don't profile the marked routine.  This flag is deprecated.
#define NOPROF	1
// It is ok for the linker to get multiple of these symbols.  It will
// pick one of the duplicates to use.
#define DUPOK	2
// Don't insert stack check preamble.
#define NOSPLIT	4
// Put this data in a read-only section.
#define RODATA	8
// This data contains no pointers.
#define NOPTR	16
// This is a wrapper function and should not count as disabling 'recover'.
#define WRAPPER 32
// end textflag.h

#define JITSTACKSIZE 2048

TEXT ·call(SB),0,$JITSTACKSIZE-24
        MOVQ b+0(FP), AX
        CALL AX
        RET

TEXT ·funcImpl(SB),0,$JITSTACKSIZE-8
        MOVQ 8(DX), AX
        MOVQ arg+0(FP), DI
        CALL AX
        RET
