package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var fileWithWords *string = flag.String("f", "./wordsEn.txt", "Path to the file that contains ALL the words")
var startWord *string = flag.String("s", "", "The word that starts the chain")
var destinationWord *string = flag.String("d", "", "The word that is the destination of the chain")
var isVerbose *bool = flag.Bool("v", false, "Verbose output")
var isDebug *bool = flag.Bool("vv", false, "Debugging output")

type WordWithPath struct {
	word string
	path []string
}

func main() {

	// Parse the command line options
	flag.Parse()

	Verbose("isVerbose:", *isVerbose)
	Verbose("isDebug:", *isDebug)
	Verbose("fileWithWords:", *fileWithWords)
	Verbose("startWord:", *startWord)
	Verbose("destinationWord:", *destinationWord)

	Verbose("##### Loading words from file " + *fileWithWords)
	words := LoadWordListFromFileAndCheckForWords(*fileWithWords, *startWord, *destinationWord)
	Debug(words)

	Verbose("##### Creating prepared lists from huge word list")
	wordMap := SplitWordListIntoSubsets(words)
	Debug(wordMap)

	Verbose("##### Searching path from '", *startWord, "' to word '", *destinationWord, "'...")
	startWordWithPath := WordWithPath{*startWord, []string{}}
	frontier := []WordWithPath{startWordWithPath}
	explored := make([]WordWithPath, 0)
	FindPath(&wordMap, frontier, explored, *destinationWord)
}

func LoadWordListFromFileAndCheckForWords(filename string, startWord string, destinationWord string) (words []string) {

	var foundStartWord bool = false
	var foundEndWord bool = false

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if word == startWord {
			foundStartWord = true
		}
		if word == destinationWord {
			foundEndWord = true
		}
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	Verbose("Loaded ", len(words), " words!")

	if !foundStartWord || !foundEndWord {
		Error("D'oh! Could not find start- or destinationWord in word list!11!!")
	}

	return
}

func SplitWordListIntoSubsets(words []string) (wordMap map[string][]string) {

	wordMap = make(map[string][]string, 0)

	for _, word := range words {
		firstTwoLetters := GetFirstTwoLetters(word)
		lastTwoLetters := GetLastTwoLetters(word)

		Verbose("Word=", word, " - ", firstTwoLetters, "...", lastTwoLetters)

		wordMap["F"+firstTwoLetters] = append(wordMap["F"+firstTwoLetters], word)
		wordMap["L"+lastTwoLetters] = append(wordMap["L"+lastTwoLetters], word)
	}
	return
}

func GetFirstTwoLetters(word string) (letters string) {
	if len(word) >= 2 {
		letters = word[0:2]
	} else {
		letters = word
	}
	return
}

func GetLastTwoLetters(word string) (letters string) {
	if len(word) >= 2 {
		letters = word[len(word)-2:]
	} else {
		letters = word
	}
	return
}

func FindPath(wordMap *map[string][]string, frontier []WordWithPath, explored []WordWithPath, destinationWord string) bool {

	for len(frontier) > 0 {

		currentWord := frontier[0]
		frontier = frontier[1:]

		if ListIncludesWord(explored, currentWord.word) {
			Debug("Word", currentWord, "already in explored... Next!", currentWord)
			continue
		}

		explored = append(explored, currentWord)

		Debug("FRONTIERSIZE: ", len(frontier), "EXPLOREDSIZE: ", len(explored))

		firstTwoLetters := GetFirstTwoLetters(currentWord.word)
		Debug("firstTwoLetters of", currentWord.word, "=", firstTwoLetters)

		lastTwoLetters := GetLastTwoLetters(currentWord.word)
		Debug("lastTwoLetters of", currentWord.word, "=", lastTwoLetters)

		possibleWords := GetPossibleWords(wordMap, "L"+firstTwoLetters, "F"+lastTwoLetters)
		Debug("POSSIBLEWORDS: ", len(possibleWords))

		for _, word := range possibleWords {

			if word == destinationWord {
				Debug("FOUND! ", word, currentWord)
				PrintWordChain(word, currentWord)
				return true
			}

			Debug("Adding", word, "to the frontier...")
			newWord := WordWithPath{}
			newWord.word = word
			newWord.path = append(currentWord.path, currentWord.word)

			frontier = append(frontier, newWord)
		}

	}

	return false
}

func PrintWordChain(lastWord string, parentWord WordWithPath) {
	Debug(len(parentWord.path))

	for _, value := range parentWord.path {
		Print(value)
	}
	Print(parentWord.word)
	Print(lastWord)
	Print("Chainsize:", 2+len(parentWord.path))
}

func ListIncludesWord(list []WordWithPath, word string) bool {
	for _, w := range list {
		if w.word == word {
			return true
		}
	}
	return false
}

func GetPossibleWords(wordMap *map[string][]string, firstTwoLetters string, lastTwoLetters string) (words []string) {
	return append((*wordMap)[firstTwoLetters], (*wordMap)[lastTwoLetters]...)
}

func Print(args ...interface{}) {
	log.Print(args)
}

func Verbose(args ...interface{}) {
	if *isVerbose || *isDebug {
		log.Print(args)
	}
}

func Debug(args ...interface{}) {
	if *isDebug {
		log.Printf("# %v", args)
	}
}

func Error(args ...interface{}) {
	log.Fatalf("ERROR %v", args)
}
