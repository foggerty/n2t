{-
Known to work with:
The Glorious Glasgow Haskell Compilation System, version 7.4.2
Extra library dependency:
Safe
To compile:
ghc -o HackAssembler Hack.hs
To run:
HackAssembler foo.asm
Will output:
foo.asm.hack
-}

module Main where
--module Text.Assembler.Hack where

import Text.ParserCombinators.Parsec
import Data.Char (isSpace)
import Data.List (isPrefixOf, unfoldr, (\\), nub, groupBy)
import Safe (readMay)
import Control.Applicative ((<$>), (<*>), (<*), (*>))
import System.Environment (getArgs)

type Label = String

data AST
     = AstI Instruction
     | AstL Label
     deriving Show

data Instruction 
     = AInstI Int  
     | AInstL Label
     | CInst { dest ::  Destination, comp :: Compute, jmp :: Jump }
     deriving Show

data Destination
     = D_Null
     | D_M
     | D_D
     | D_MD
     | D_A
     | D_AM
     | D_AD
     | D_AMD
     deriving (Show, Read)

data Compute 
     = C_Zero 
     | C_One 
     | C_NegOne 
     | C_D 
     | C_A
     | C_M
     | C_NotD
     | C_NotA
     | C_NotM
     | C_NegD
     | C_NegA
     | C_NegM
     | C_DPlusOne
     | C_APlusOne
     | C_MPlusOne
     | C_DMinOne
     | C_AMinOne
     | C_MMinOne
     | C_DPlusA
     | C_DPlusM
     | C_DMinA
     | C_DMinM
     | C_AMinD
     | C_MMinD
     | C_DAndA
     | C_DAndM
     | C_DOrA
     | C_DOrM
     deriving Show

data Jump
     = J_NULL
     | JGT
     | JEQ
     | JGE
     | JLT
     | JNE
     | JLE
     | JMP
     deriving (Show, Read)

run p input =  case (parse p "" input) of
  Left err -> putStr "parse error at" >> print err
  Right x  -> print x

filterOut :: (a -> Bool) -> [a] -> [a]
filterOut f xs = filter (not . f) xs

purgeInput :: [String] -> [String]
purgeInput = filterOut badLine . map dehydrate
    where dehydrate = filterOut isSpace
          badLine l = isPrefixOf "//" l || length l == 0

parseFileContents :: String -> [AST]
parseFileContents =  map parseLine . purgeInput . lines

parseLine :: String -> AST
parseLine s = case (parse parseAst "" s) of
                Left err -> error $ show err
                Right x -> x

parseLabel :: Parser Label
parseLabel =
    do f <- letter <|> oneOf "_.$:"
       rest <- option "" $ many1 $ letter <|> digit <|> oneOf "_.$:"
       return (f:rest)

parseAst :: Parser AST
parseAst =
  do p <- pL <|> pI
     optional (string "//" >> optional (many1 anyChar))
     eof
     return p
  where pL = fmap AstL $ between (char '(') (char ')') parseLabel
        pI = fmap AstI parseInst

parseInst :: Parser Instruction
parseInst = parseAInst <|> parseCInst

parseAInst :: Parser Instruction
parseAInst = 
    do char '@'
       p <- parseAInstI <|> parseAInstL
       return p

parseAInstI :: Parser Instruction
parseAInstI =
    do d <- many1 digit
       (return . AInstI . read) d

parseAInstL :: Parser Instruction
parseAInstL =
    do p <- parseLabel
       (return . AInstL) p

parseCInst :: Parser Instruction
parseCInst =
    do d <- parseDest
       c <- parseCmd
       j <- parseJmp
       return $ CInst d c j

parseDest :: Parser Destination
parseDest = option D_Null (try parse)
  where parse = do s <- many1 $ oneOf "MDA"
                   d <- case (readMay ("D_" ++ s) :: Maybe Destination) of
                     Just x -> return x
                     Nothing -> error "should not reach this is parseDest"
                   char '='
                   return d

cmdLookup :: [(String, Compute)]
cmdLookup = [("0", C_Zero)
            ,("1", C_One)
            ,("-1", C_NegOne)
            ,("D", C_D)
            ,("A", C_A)
            ,("M", C_M)
            ,("!D", C_NotD)
            ,("!A", C_NotA)
            ,("!M", C_NotM)
            ,("-D", C_NegD)
            ,("-A", C_NegA)
            ,("-M", C_NegM)
            ,("D+1", C_DPlusOne)
            ,("A+1", C_APlusOne)
            ,("M+1", C_MPlusOne)
            ,("D-1", C_DMinOne)
            ,("A-1", C_AMinOne)
            ,("M-1", C_MMinOne)
            ,("D+A", C_DPlusA)
            ,("D+M", C_DPlusM)
            ,("D-A", C_DMinA)
            ,("D-M", C_DMinM)
            ,("A-D", C_AMinD)
            ,("M-D", C_MMinD)
            ,("D&A", C_DAndA)
            ,("D&M", C_DAndM)
            ,("D|A", C_DOrA)
            ,("D|M", C_DOrM)
            ]

parseCmd :: Parser Compute
parseCmd = 
  do s <- many1 $ oneOf "01MDA+-!&|"
     case lookup s cmdLookup of
       Just x -> return x
       Nothing -> error "should not reach this in parseCmd"

parseJmp :: Parser Jump
parseJmp =
  parse <|> (return J_NULL)
  where parse = do char ';'
                   s <- many1 $ oneOf "JMPEQLTGN"
                   case (readMay s :: Maybe Jump) of
                     Just x -> return x
                     Nothing -> error "should not reach this is parseJmp"

toHackAsm :: [(Label,Int)] -> Instruction -> String
toHackAsm _ (AInstI n) = pad 16 '0' $ i2b n
toHackAsm jt (AInstL l) = pad 16 '0' $ i2b n
  where n = case lookup l jt of
                 Just x -> x
                 Nothing -> error "bad lookup"
toHackAsm _ (CInst d c j) = pad 16 '1' $ cs ++ ds ++ js
  where ds = toHackAsmD d
        cs = toHackAsmC c
        js = toHackAsmJ j

toHackAsmD :: Destination -> String
toHackAsmD D_Null = "000" 
toHackAsmD D_M    = "001"
toHackAsmD D_D    = "010"
toHackAsmD D_MD   = "011"
toHackAsmD D_A    = "100"
toHackAsmD D_AM   = "101"
toHackAsmD D_AD   = "110"
toHackAsmD D_AMD  = "111"

toHackAsmC :: Compute -> String
toHackAsmC C_Zero     = "0101010"
toHackAsmC C_One      = "0111111"
toHackAsmC C_NegOne   = "0111010"
toHackAsmC C_D        = "0001100"
toHackAsmC C_A        = "0110000"
toHackAsmC C_M        = "1110000"
toHackAsmC C_NotD     = "0001101"
toHackAsmC C_NotA     = "0110001"
toHackAsmC C_NotM     = "1110001"
toHackAsmC C_NegD     = "0001111"
toHackAsmC C_NegA     = "0110011"
toHackAsmC C_NegM     = "1110011"
toHackAsmC C_DPlusOne = "0011111"
toHackAsmC C_APlusOne = "0110111"
toHackAsmC C_MPlusOne = "1110111"
toHackAsmC C_DMinOne  = "0001110"
toHackAsmC C_AMinOne  = "0110010"
toHackAsmC C_MMinOne  = "1110010"
toHackAsmC C_DPlusA   = "0000010"
toHackAsmC C_DPlusM   = "1000010"
toHackAsmC C_DMinA    = "0010011"
toHackAsmC C_DMinM    = "1010011"
toHackAsmC C_AMinD    = "0000111"
toHackAsmC C_MMinD    = "1000111"
toHackAsmC C_DAndA    = "0000000"
toHackAsmC C_DAndM    = "1000000"
toHackAsmC C_DOrA     = "0010101"
toHackAsmC C_DOrM     = "1010101"

toHackAsmJ :: Jump -> String
toHackAsmJ J_NULL = "000"
toHackAsmJ JGT    = "001"
toHackAsmJ JEQ    = "010"
toHackAsmJ JGE    = "011"
toHackAsmJ JLT    = "100"
toHackAsmJ JNE    = "101"
toHackAsmJ JLE    = "110"
toHackAsmJ JMP    = "111"

i2b :: Int -> String
i2b n = reverse $ unfoldr f n
    where f x
            | x > 0 = Just (qr x) 
            | otherwise = Nothing

qr :: Int -> (Char, Int)
qr n = (binary r, q)
    where (q,r) = quotRem n 2
          binary 0 = '0'
          binary 1 = '1'
          binary _ = error "Not a binary digit"

pad :: Int -> a -> [a] -> [a]
pad n a xs = (take y (repeat a)) ++ xs
    where y = n - (length xs)

reservedLabels :: [(Label, Int)]
reservedLabels = [ ("SP", 0)
                 , ("LCL", 1)
                 , ("ARG", 2)
                 , ("THIS", 3)
                 , ("THAT", 4)
                 , ("SCREEN", 16384)
                 , ("KBD", 24576)
                 ] ++ [("R" ++ show x, x) | x <- [0..15]]

makeLabelTbl :: [AST] -> [(Label, Int)]
makeLabelTbl = fold 0 []
  where
    fold nexti tbl [] = tbl
    fold nexti tbl ((AstI _):xs) = fold (nexti + 1) tbl xs
    fold nexti tbl ((AstL l):xs) = fold nexti ((l,nexti):tbl) xs

filterInst :: [AST] -> [Instruction]
filterInst = map m . filter f
  where
    f (AstI _) = True
    f (AstL _) = False
    m (AstI i) = i
    
filterAInstL :: [Instruction] -> [Instruction]
filterAInstL = filter f
  where
    f (AInstL _) = True
    f _          = False

assemble :: String -> String
assemble prg = unlines $ map (toHackAsm (jt++st)) instructions
  where
    p = parseFileContents prg
    instructions = filterInst p
    labeledInst = filterAInstL instructions
    symbols = (nub (map (\(AInstL l) -> l) labeledInst)) \\ (map fst jt)
    st = zip symbols [16..]
    jt = (makeLabelTbl p) ++ reservedLabels

main = do
  file <- head <$> getArgs
  source <- readFile file
  writeFile (file ++ ".hack") (assemble source)