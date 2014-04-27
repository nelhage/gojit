package amd64

const (
	PREFIX_LOCK     = 0xF0
	PREFIX_REPNZ    = 0xF2
	PREFIX_REPZ     = 0xF3
	PREFIX_SEG_CS   = 0x2E
	PREFIX_SEG_SS   = 0x36
	PREFIX_SEG_DS   = 0x3E
	PREFIX_SEG_ES   = 0x26
	PREFIX_SEG_FS   = 0x64
	PREFIX_SEG_GS   = 0x65
	PREFIX_OPSIZE   = 0x66
	PREFIX_ADDRSIZE = 0x67

	MOD_INDIR        = 0x0
	MOD_INDIR_DISP8  = 0x1
	MOD_INDIR_DISP32 = 0x2
	MOD_REG          = 0x3

	SCALE_1 = 0x0
	SCALE_2 = 0x1
	SCALE_4 = 0x2
	SCALE_8 = 0x3

	/* overflow */
	CC_O  = 0x0
	CC_NO = 0x1
	/* unsigned comparisons */
	CC_B  = 0x2
	CC_AE = 0x3
	CC_BE = 0x6
	CC_A  = 0x7
	/* zero */
	CC_Z  = 0x4
	CC_NZ = 0x5
	/* sign */
	CC_S  = 0x8
	CC_NS = 0x9
	/* parity */
	CC_P  = 0xA
	CC_NP = 0xB
	/* unsigned comparisons */
	CC_L  = 0xC
	CC_GE = 0xD
	CC_LE = 0xE
	CC_G  = 0xF

	PFX_REX = 0x40
	REXW    = 0x08
	REXR    = 0x04
	REXX    = 0x02
	REXB    = 0x01

	REG_DISP32 = 5
	REG_SIB    = 4
)
