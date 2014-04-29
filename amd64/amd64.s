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


TEXT ·get_runtime_cgocallback_gofunc(SB),0,$0-8
        LEAQ runtime·cgocallback_gofunc(SB), AX
        MOVQ AX, rv+0(FP)
        RET
