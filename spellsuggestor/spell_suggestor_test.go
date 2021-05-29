package spellsuggestor

import (
	"reflect"
	"testing"
)

func TestSpellSuggestor_Suggest(t *testing.T) {
	spellSuggestor, err := New("../dictionary.txt")
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		dictionary map[string]bool
	}
	type args struct {
		word string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{"abbtos -> abbots", fields{spellSuggestor.dictionary}, args{"abbtos"}, []string{"abbots"}},
		{"aebreviator -> abbreviator", fields{spellSuggestor.dictionary}, args{"aebreviator"}, []string{"abbreviator"}},
		{"aerreviator -> abbreviator", fields{spellSuggestor.dictionary}, args{"aerreviator"}, []string{"abbreviator"}},
		{"atylor -> tailor", fields{spellSuggestor.dictionary}, args{"atylor"}, []string{"tailor"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SpellSuggestor := &SpellSuggestor{
				dictionary: tt.fields.dictionary,
			}
			if got := SpellSuggestor.Suggest(tt.args.word); !gotWhatWeWant(got, tt.want) {
				t.Errorf("Suggest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func gotWhatWeWant(got, want []string) bool {
	gotMap := map[string]bool{}
	for i := range got {
		gotMap[got[i]] = true
	}
	for i := range want {
		if !gotMap[want[i]] {
			return false
		}
	}
	return true
}

func TestSpellSuggestor_getWordsOneEditAway(t *testing.T) {
	expectedSuggestions := []string{
		// add
		"aab", "bab", "cab", "dab", "eab", "fab", "gab", "hab", "iab", "jab", "kab", "lab", "mab", "nab", "oab", "pab", "qab", "rab", "sab", "tab", "uab", "vab", "wab", "xab", "yab", "zab",
		// swap
		// replace
		"bb", "cb", "db", "eb", "fb", "gb", "hb", "ib", "jb", "kb", "lb", "mb", "nb", "ob", "pb", "qb", "rb", "sb", "tb", "ub", "vb", "wb", "xb", "yb", "zb",
		// delete
		"b",

		// add
		"aab", "abb", "acb", "adb", "aeb", "afb", "agb", "ahb", "aib", "ajb", "akb", "alb", "amb", "anb", "aob", "apb", "aqb", "arb", "asb", "atb", "aub", "avb", "awb", "axb", "ayb", "azb",
		// swap
		"ba",
		// replace
		"aa", "ac", "ad", "ae", "af", "ag", "ah", "ai", "aj", "ak", "al", "am", "an", "ao", "ap", "aq", "ar", "as", "at", "au", "av", "aw", "ax", "ay", "az",

		// delete
		"a",

		// add at end
		"aba", "abb", "abc", "abd", "abe", "abf", "abg", "abh", "abi", "abj", "abk", "abl", "abm", "abn", "abo", "abp", "abq", "abr", "abs", "abt", "abu", "abv", "abw", "abx", "aby", "abz",
	}

	type fields struct {
		dictionary map[string]bool
	}
	type args struct {
		word string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{"simple test", fields{make(map[string]bool)}, args{"ab"}, expectedSuggestions},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spellSuggestor := &SpellSuggestor{
				dictionary: tt.fields.dictionary,
			}
			if got := spellSuggestor.getWordsOneEditAway(tt.args.word); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWordsOneEditAway() = %v, want %v", got, tt.want)
			}
		})
	}
}

func haveSameWords(a, b []string) bool {
	if len(a) == len(b) {
		hashmapCount := map[string]int{}
		for i := range a {
			hashmapCount[a[i]]++
			hashmapCount[b[i]]--
		}
		for _, count := range hashmapCount {
			if count != 0 {
				return false
			}
		}
	}
	return false
}

func unique(words []string) []string {
	hashmap := map[string]bool{}
	result := make([]string, len(words))
	i := 0
	for _, word := range words {
		if !hashmap[word] {
			result[i] = word
			i++
		}
	}
	return result[:i]
}
