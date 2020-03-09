package moveGeneration

import (
	"CheckerServer/server/common/stack"
	"CheckerServer/server/common/util"
	"fmt"
	"math"
	"time"
)

// squareToBitboard确定paddedArray中squareIndex的棋子对应Bitboard中的索引
func squareToBitboard(squareIndex int8/*padded array中的索引号*/) int8 {
	// squareIndex不在paddedArray范围内
	if squareIndex < 6 || squareIndex > 60 || squareIndex % 11 == 0 {
		return -1
	}
	// squareIndex是paddedArray的棋子
	return int8(bitboardArray[squareIndex])
}

// bitboardToSquare确定bitboard中bitboardIndex的棋子对应padded array中的索引
func bitboardToSquare(bitboardIndex int8) int8{
	// bitboardIndex不在bitboard范围内
	if bitboardIndex < 0 || bitboardIndex > 49 {
		return -1
	}
	// bitboardIndex是bitboard的棋子
	return int8(paddedArray[bitboardIndex])
}

// setpaddedArrayBoard根据bitboardPiece设置paddedArrayBoard的piece棋子
func setpaddedArrayBoard (piece int8/*WPIECE/BPIECE/WKING/BKING*/, bitboardPiece uint64) {
	for bitboardPiece > 0 {
		index := uint(util.FindLowBit(bitboardPiece))
		bitboardPiece &= ^(1 << index)
		paddedArrayBoard[paddedArray[index]] = piece
	}
}

// bitboardToPaddedArray建立一个以padded array的索引布局的棋盘
func BitboardToPaddedArrayBoard(cb *CheckerBitboard) {
	temp := cb.W
	setpaddedArrayBoard(WPIECE, temp)
	temp = cb.B
	setpaddedArrayBoard(BPIECE, temp)
	temp = cb.W & cb.K
	setpaddedArrayBoard(WKING, temp)
	temp = cb.B & cb.K
	setpaddedArrayBoard(BKING, temp)
}

////////////////////////////////////////////////////////
// checkerBitboard start
////////////////////////////////////////////////////////
var maskL4 uint64
var maskL6 uint64
var maskR4 uint64
var maskR6 uint64
var paddedArray [50]uint8
var bitboardArray [61]uint8
var paddedArrayBoard [61]int8
var counter [61]int8
const INVALID int8 = -99
const EMPTY int8 = 0
const WPIECE int8 = 1
const BPIECE int8 = 2
const WKING int8 = -1
const BKING int8 = -2

type CheckerBitboard struct {
	W,B,K uint64
	nOccCut uint64
}

// kingJumpSearchDirectionWhite在指定方向DIR上查看白王是否有跳吃可能，如可能，则将该王棋加入movers中
func (cb *CheckerBitboard) kingJumpSearchDirectionWhite(movers *uint64/*可跳吃的棋子*/,squareKingIndex uint8/*王棋在padded array中的索引*/, DIR int8/*棋子移动方向：+5，+6，-5，-6*/, index uint/*王棋在bitboard中的索引*/, nOccCut uint64/*bitboard中的空棋位*/, flag *int/*标识此方向是否能够跳吃*/) {
	for squareIndex := int8(squareKingIndex) + DIR; squareToBitboard(squareIndex + 2*DIR) >= 0 ; squareIndex += DIR {
		if 1 << uint(squareToBitboard(squareIndex)) & cb.W > 0 {
			break
		}
		if (1 << uint(squareToBitboard(squareIndex+DIR)) & cb.B) > 0 && (1 << uint(squareToBitboard(squareIndex+2*DIR)) & nOccCut) > 0 {
			*movers |= 1 << uint(index)
			*flag = 1
			break
		}
	}
}

// kingJumpSearchDirectionBlack在指定方向DIR上查看黑王是否有跳吃可能，如可能，则将该王棋加入movers中
func (cb *CheckerBitboard) kingJumpSearchDirectionBlack(movers *uint64/*可跳吃的棋子*/,squareKingIndex uint8/*王棋在padded array中的索引*/, DIR int8/*棋子移动方向：+5，+6，-5，-6*/, index uint/*王棋在bitboard中的索引*/, nOccCut uint64/*bitboard中的空棋位*/, flag *int/*标识此方向是否能够跳吃*/) {
	for squareIndex := int8(squareKingIndex) + DIR; squareToBitboard(squareIndex + 2*DIR) >= 0 ; squareIndex += DIR {
		if 1 << uint(squareToBitboard(squareIndex)) & cb.B > 0 {
			break
		}
		if (1 << uint(squareToBitboard(squareIndex+DIR)) & cb.W) > 0 && (1 << uint(squareToBitboard(squareIndex+2*DIR)) & nOccCut) > 0 {
			*movers |= 1 << uint(index)
			*flag = 1
			break
		}
	}
}

// getMoversWhite用于求出可移动的白子棋位
func (cb *CheckerBitboard) GetMoversWhite( ) uint64 {
	nOcc := uint64(^(cb.B | cb.W))
	nOccCut := util.CutOver49(nOcc)
	// 获取可往前移动一步的白子位置
	moves := nOccCut >> 5
	moves |= nOccCut&maskR4 >> 4
	moves |= nOccCut&maskR6 >> 6
	moves &= cb.W
	// 获取可移动的白王位置
	WK := cb.W & cb.K
	if WK > 0 { /* 有必要进行这个判断吗？如果没有白王，后续结果也不会影响，因为需要进行且运算，下同 */
		moves |= (nOccCut << 5) & WK
		moves |= (nOccCut&maskL4 << 4) & WK
		moves |= (nOccCut&maskL6 << 6) & WK
	}
	movesCut := util.CutOver49(moves)
	return movesCut
}

// getMoversBlack用于求出可移动的黑子棋位
func (cb *CheckerBitboard) GetMoversBlack( ) uint64 {
	nOcc := uint64(^(cb.B | cb.W))
	nOccCut := util.CutOver49(nOcc)
	// 获取可往前移动一步的黑子位置
	moves := nOccCut << 5
	moves |= nOccCut&maskL4 << 4
	moves |= nOccCut&maskL6 << 6
	moves &= cb.B
	// 获取可移动的黑王位置
	BK := cb.B & cb.K
	if BK > 0 {
		moves |= (nOccCut >> 5) & BK
		moves |= (nOccCut&maskR4 >> 4) & BK
		moves |= (nOccCut&maskR6 >> 6) & BK
	}
	movesCut := util.CutOver49(moves)
	return movesCut
}

// getJumpersWhite用于求出可以跳吃黑子的白子棋位
func (cb *CheckerBitboard) GetJumpersWhite( ) uint64 {
	nOcc := uint64(^(cb.B | cb.W))
	nOccCut := util.CutOver49(nOcc)
	// 获取可以跳吃黑子的白子棋位-往前单步跳吃
	movers := uint64(0)
	temp := (nOccCut >> 5) & cb.B
	if temp > 0 { movers |= (((temp & maskR4) >> 4) | ((temp & maskR6) >> 6)) & cb.W }
	temp = (((nOccCut & maskR4) >> 4) | ((nOccCut & maskR6) >> 6)) & cb.B
	if temp > 0 { movers |= (temp >> 5) & cb.W }
	// 获取可以跳吃黑子的白子棋位-往后单步跳吃
	temp = (nOccCut << 5) & cb.B
	if temp > 0 { movers |= (((temp & maskL4) << 4) | ((temp & maskL6) << 6)) & cb.W }
	temp = (((nOccCut & maskL4) << 4) | ((nOccCut & maskL6) << 6)) & cb.B
	if temp > 0 { movers |= (temp << 5) & cb.W }

	// 获取可以跳吃黑子的白王棋位
	WK := cb.W & cb.K
	WK &= ^movers // 已在movers中的白王不必重复计算
	for WK > 0 {
		index := uint(util.FindLowBit(WK))
		WK &= ^(1 << index)
		squareWK := paddedArray[index]

		flag := 0

		// 获取可以跳吃黑子的白王棋位-向前
		cb.kingJumpSearchDirectionWhite(&movers ,squareWK, 5/*+5，+6，-5，-6*/, index, nOccCut, &flag )
		if flag == 1 { break }
		cb.kingJumpSearchDirectionWhite(&movers ,squareWK, 6/*+5，+6，-5，-6*/, index, nOccCut, &flag )
		if flag == 1 { break }

		// 获取可以跳吃黑子的白王棋位-向后
		cb.kingJumpSearchDirectionWhite(&movers ,squareWK, -5/*+5，+6，-5，-6*/, index, nOccCut, &flag )
		if flag == 1 { break }
		cb.kingJumpSearchDirectionWhite(&movers ,squareWK, -6/*+5，+6，-5，-6*/, index, nOccCut, &flag )
	}
	return movers
}

// getJumpersBlack用于求出可以跳吃白子的黑子棋位
func (cb *CheckerBitboard) GetJumpersBlack( ) uint64 {
	nOcc := uint64(^(cb.B | cb.W))
	nOccCut := util.CutOver49(nOcc)
	// 获取可以跳吃白子的黑子棋位-往前单步跳吃
	movers := uint64(0)
	temp := (nOccCut << 5) & cb.W
	if temp > 0 { movers |= (((temp & maskL4) << 4) | ((temp & maskL6) << 6)) & cb.B }
	temp = (((nOccCut & maskL4) << 4) | ((nOccCut & maskL6) << 6)) & cb.W
	if temp > 0 { movers |= (temp << 5) & cb.B }
	// 获取可以跳吃白子的黑子棋位-往后单步跳吃
	temp = (nOccCut >> 5) & cb.W
	if temp > 0 { movers |= (((temp & maskR4) >> 4) | ((temp & maskR6) >> 6)) & cb.B }
	temp = (((nOccCut & maskR4) >> 4) | ((nOccCut & maskR6) >> 6)) & cb.W
	if temp > 0 { movers |= (temp >> 5) & cb.B }

	// 获取可以跳吃白子的黑王棋位
	BK := cb.B & cb.K
	BK &= ^movers // 已在movers中的黑王不必重复计算
	for BK > 0 {
		index := uint(util.FindLowBit(BK))
		BK &= ^(1 << index)
		squareBK := paddedArray[index]

		flag := 0

		// 获取可以跳吃白子的黑王棋位-向前
		cb.kingJumpSearchDirectionBlack(&movers ,squareBK, -5/*+5，+6，-5，-6*/, index, nOccCut, &flag)
		if flag == 1 { break }
		cb.kingJumpSearchDirectionBlack(&movers ,squareBK, -6/*+5，+6，-5，-6*/, index, nOccCut, &flag)
		if flag == 1 { break }

		// 获取可以跳吃白子的黑王棋位-向后
		cb.kingJumpSearchDirectionBlack(&movers ,squareBK, 5/*+5，+6，-5，-6*/, index, nOccCut, &flag)
		if flag == 1 { break }
		cb.kingJumpSearchDirectionBlack(&movers ,squareBK, 6/*+5，+6，-5，-6*/, index, nOccCut, &flag)
	}
	return movers
}

// getBlackMoves用于求出可移动的黑子个数（不含跳吃）
func (cb *CheckerBitboard) getBlackMoves( ) int {
	numBlackMove := util.CountUint64(cb.GetMoversBlack())
	return numBlackMove
}

// getWhiteMoves用于求出可移动的白子个数（不含跳吃）
func (cb *CheckerBitboard) getWhiteMoves( ) int {
	numWhiteMove := util.CountUint64(cb.GetMoversWhite())
	return numWhiteMove
}

////////////////////////////////////////////////////////
// checkerBitboard end
////////////////////////////////////////////////////////


////////////////////////////////////////////////////////
// checkerMoveList start
////////////////////////////////////////////////////////
type checkerMove struct {
	parent *checkerMove	// 父节点
	src,dst uint8
	layer uint8	// 该节点属于第几层
}

type CheckerMoveList struct {
	numMoves, numJumps uint8 // 需要吗？
	moves[100] checkerMove // 大小设置有待考量。对于普通移动，存移动本身；对于跳吃，存叶子节点。
}

// ownPiece根据对手棋子的颜色oppoPiece求得本方棋子的颜色
func ownPiece(oppoPiece uint8) uint8 {
	return -oppoPiece+3
}

// isIdxValid确认idx在padded array中是否合法
func isIdxValid(idx int8) bool {
	if idx >= 6 && idx <= 60 {
		if paddedArrayBoard[idx] != INVALID {
			return true
		} else {
			return false
		}
	}
	return false
}

// 清空计数器
func resetCounter() {
	counter = [61]int8{INVALID, INVALID, INVALID, INVALID, INVALID, INVALID,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0}
}

// clear把成员numMoves和numJumps赋零
func (cml *CheckerMoveList) clear() {
	cml.numMoves = 0
	cml.numJumps = 0
}

// createMove创建跳吃的节点
func (cml *CheckerMoveList) createMove(parent *checkerMove,src uint8, dst uint8, layer uint8) *checkerMove{
	move := new(checkerMove)
	move.parent = parent
	move.src = src
	move.dst = dst
	move.layer = layer
	return move
}

// addMove给moves添加普通移动
func (cml *CheckerMoveList) addMove(src uint8, dst uint8) {
	move := cml.createMove(nil,src,dst,0)
	cml.moves[cml.numMoves] = *move
	cml.numMoves ++
}

// findMovesWhite找出当前bitboard上白子的所有可走步，并将其存放于checkerMoveList中
func (cml *CheckerMoveList) findMovesWhite(cb *CheckerBitboard) {
	jumpers := cb.GetJumpersWhite()
	if jumpers > 0 {
		cml.findJumps(jumpers, uint8(BPIECE))
	} else { // 如果能跳吃，就无需寻找可移动的棋子，因为有跳吃必须走跳吃
		movers := cb.GetMoversWhite()
		cml.findNonJumps(movers, uint8(BPIECE))
	}
}

// findMovesBlack找出当前bitboard上黑子的所有可走步，并将其存放于checkerMoveList中
func (cml *CheckerMoveList) findMovesBlack(cb *CheckerBitboard) {
	jumpers := cb.GetJumpersBlack()
	if jumpers > 0 {
		cml.findJumps(jumpers, uint8(WPIECE))
	} else { // 如果能跳吃，就无需寻找可移动的棋子，因为有跳吃必须走跳吃
		movers := cb.GetMoversBlack()
		cml.findNonJumps(movers, uint8(WPIECE))
	}
}

// findJumps找出当前bitboard上白/黑子的所有可跳吃（包括连跳）步，并将其存放于checkerMoveList中
func (cml *CheckerMoveList) findJumps(jumpers uint64, opponentPiece uint8 /*对手棋子*/) {
	cml.clear()
	for jumpers > 0 {
		index := uint8(util.FindLowBit(jumpers))
		jumpers &= ^(1 << index)
		square := bitboardToSquare(int8(index))

		numSqMoves := int8(0)
		layer := uint8(0)
		// 四个方向寻找跳吃
		if square >= 0 { // 根节点没有父节点
			cml.checkJumpDir(nil, uint8(square), 5, opponentPiece, layer,&numSqMoves)
			cml.checkJumpDir(nil, uint8(square), 6, opponentPiece, layer,&numSqMoves)
			cml.checkJumpDir(nil, uint8(square), -5, opponentPiece, layer,&numSqMoves)
			cml.checkJumpDir(nil, uint8(square), -6, opponentPiece, layer,&numSqMoves)
		} else {
			fmt.Printf("index不在bitboard范围内！")
		}
	}
	cml.numMoves = cml.numJumps
}

// findJumps找出当前bitboard上白/黑子的所有可移动步，并将其存放于checkerMoveList中
func (cml *CheckerMoveList) findNonJumps(movers uint64, opponentPiece uint8 /*对手棋子*/) {
	cml.clear()
	for movers > 0 {
		index := uint8(util.FindLowBit(movers))
		movers &= ^(1 << index)
		square := bitboardToSquare(int8(index))

		if paddedArrayBoard[square] == BKING || paddedArrayBoard[square] == WKING { // 王棋四个方向都能走
			cml.checkNonJumpDir(uint8(square), 5)
			cml.checkNonJumpDir(uint8(square), -5)
			cml.checkNonJumpDir(uint8(square), 6)
			cml.checkNonJumpDir(uint8(square), -6)
		} else {
			if opponentPiece == uint8(BPIECE) { // 白子只能往上走
				if paddedArrayBoard[square + 5] == EMPTY { cml.addMove(uint8(square),uint8(square+5)) }
				if paddedArrayBoard[square + 6] == EMPTY { cml.addMove(uint8(square),uint8(square+6))}
			}
			if opponentPiece == uint8(WPIECE) { // 黑子只能往下走
				if paddedArrayBoard[square - 5] == EMPTY { cml.addMove(uint8(square),uint8(square-5)) }
				if paddedArrayBoard[square - 6] == EMPTY { cml.addMove(uint8(square),uint8(square-6))}
			}
		}
	}
}

// checkJumpDir确定DIR有无跳吃，如果有，继续寻找连跳
func (cml *CheckerMoveList) checkJumpDir(parent *checkerMove, square uint8 /*当前棋位在padded array上的索引*/, DIR int8 /*方向*/, opponentPiece uint8 /*对手棋子*/, layer uint8, numSqMoves *int8) {
	moveStack := stack.NewStack()
	var jumpedPieceSq uint8

	// 存moveStack
	if paddedArrayBoard[square] == BKING || paddedArrayBoard[square] == WKING {	// 王的跳吃
		for squareIndex := int8(square); isIdxValid(int8(squareIndex) + 2*DIR) && isIdxValid(int8(squareIndex) + DIR); squareIndex += DIR {
			if int8(math.Abs(float64(paddedArrayBoard[uint8(int8(squareIndex) + DIR)]))) == int8(ownPiece(opponentPiece)) || counter[squareIndex + DIR] != 0 {	// 下一位是本方棋子，该方向不存在跳吃可能
				break
			}
			if int8(math.Abs(float64(paddedArrayBoard[uint8(int8(squareIndex) + DIR)]))) == int8(opponentPiece) && paddedArrayBoard[uint8(int8(squareIndex) + 2*DIR)] == EMPTY && counter[squareIndex + 2 * DIR] == 0{	// 下一位是对手棋子且再下一位为空，可以跳吃
				*numSqMoves ++
				jumpedPieceSq = uint8(int8(squareIndex) + DIR)
				for idx := int8(squareIndex) + 2*DIR; isIdxValid(idx) && counter[idx] == 0; idx += DIR {
					if paddedArrayBoard[idx] == EMPTY {
						move := cml.createMove(parent,square,uint8(idx),layer)
						moveStack.Push(move)
					} else {
						break
					}
				}
				break
			}
		}
	} else {	// 兵的跳吃
		if isIdxValid(int8(square) + 2*DIR) && isIdxValid(int8(square) + DIR) && int8(math.Abs(float64(paddedArrayBoard[uint8(int8(square) + DIR)]))) == int8(opponentPiece) && paddedArrayBoard[uint8(int8(square) + 2*DIR)] == EMPTY { // 下一位是对手棋子且再下一位为空，可以跳吃
			*numSqMoves ++
			jumpedPieceSq = uint8(int8(square) + DIR)
			move := cml.createMove(parent,square,uint8(int8(square) + 2*DIR),layer)
			moveStack.Push(move)
		}
	}

	layer = layer + 1
	for ;!moveStack.Empty(); {
		if layer == 1 {
			resetCounter()
		}
		topElement := moveStack.Top()
		move := topElement.(*checkerMove)

		// 查找后续跳吃
		cml.findSqJumps(move,layer,jumpedPieceSq, opponentPiece)

		moveStack.Pop()
	}
}

// findSqJumps是跟checkJumpDir打配合哒
func (cml *CheckerMoveList) findSqJumps(move *checkerMove, layer uint8, jumpedPieceSq uint8, opponentPiece uint8 /*对手棋子*/) {
	numSqMoves := int8(0)

	// 将被跳吃的棋子“拿走”，以防多次连跳
	jumpedPiece := paddedArrayBoard[jumpedPieceSq]
	paddedArrayBoard[jumpedPieceSq] = EMPTY
	counter[jumpedPieceSq] += 1
	jumpPiece := paddedArrayBoard[move.src]
	paddedArrayBoard[move.src] = EMPTY
	paddedArrayBoard[move.dst] = jumpPiece

	// 寻找连跳
	cml.checkJumpDir(move, move.dst, 5, opponentPiece, layer, &numSqMoves)
	cml.checkJumpDir(move, move.dst, 6, opponentPiece, layer, &numSqMoves)
	cml.checkJumpDir(move, move.dst, -5, opponentPiece, layer, &numSqMoves)
	cml.checkJumpDir(move, move.dst, -6, opponentPiece, layer, &numSqMoves)

	// 叶子节点，存
	if numSqMoves == 0 {
		cml.moves[cml.numJumps] = *move
		cml.numJumps ++
	}

	// 把被“拿走”的棋子放回，还原棋盘
	paddedArrayBoard[jumpedPieceSq] = jumpedPiece
	paddedArrayBoard[move.src] = jumpPiece
	paddedArrayBoard[move.dst] = EMPTY
	counter[jumpedPieceSq] -= 1
}

// checkJumpDir确定DIR有无王的移动
func (cml *CheckerMoveList) checkNonJumpDir(square uint8 /*当前棋位在padded array上的索引*/, DIR int8 /*方向*/ ) {
	for idx := int8(square); isIdxValid(idx + DIR); idx += DIR {
		if paddedArrayBoard[idx + DIR] != EMPTY {
			break
		} else {
			cml.addMove(uint8(square),uint8(idx+DIR))
		}
	}
}
////////////////////////////////////////////////////////
// checkerMoveList end
////////////////////////////////////////////////////////


func init() {
	maskL4 = 33017592576030	    // 向左上方行走数字+4的位置：01 02 03 04 11 12 13 14 21 22 23 24 31 32 33 34 41 42 43 44
	maskL6 = 515899884000		// 向右上方行走数字+6的位置：05 06 07 08 15 16 17 18 25 26 27 28 35 36 37 38
	maskR4 = 528281481216480	// 向右下方行走数字-4的位置：45 46 47 48 35 36 37 38 25 26 27 28 15 16 17 18 05 06 07 08
	maskR6 = 33017592576000		// 向左下方行走数字-6的位置：41 42 43 44 31 32 33 34 21 22 23 24 11 12 13 14
	// padded array中棋盘的索引集合
	paddedArray = [50]uint8{6,7,8,9,10,12,13,14,15,16,17,18,19,20,21,23,24,25,26,27,28,29,30,31,32,34,35,36,37,38,39,40,41,42,43,45,46,47,48,49,50,51,52,53,54,56,57,58,59,60}
	bitboardArray = [61]uint8{0,0,0,0,0,0,0,1,2,3,4,0,5,6,7,8,9,10,11,12,13,14,0,15,16,17,18,19,20,21,22,23,24,0,25,26,27,28,29,30,31,32,33,34,0,35,36,37,38,39,40,41,42,43,44,0,45,46,47,48,49}
	paddedArrayBoard = [61]int8{INVALID,INVALID,INVALID,INVALID,INVALID,INVALID,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY,	INVALID,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY,	INVALID,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY,	INVALID,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY,	INVALID,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY,	INVALID,
		EMPTY,	EMPTY,	EMPTY,	EMPTY,	EMPTY}
	counter = [61]int8{INVALID, INVALID, INVALID, INVALID, INVALID, INVALID,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, INVALID,
		0, 0, 0, 0, 0}
}

func main() {
	start := time.Now().UnixNano()

	// 初始化bitboard
	var cb = &CheckerBitboard{W:1048608,B:103085506884,K:4194400}
	// 初始化padded array board
	BitboardToPaddedArrayBoard(cb)

	var cml CheckerMoveList
	// 获取可以移动的白子的移动路径
	cml.findMovesWhite(cb)

	end := time.Now().UnixNano()
	timeElapsed := end - start
	//fmt.Printf("%b\n",check)
	fmt.Printf("耗时（毫秒）：%v;\n",timeElapsed)
}
