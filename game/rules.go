package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Setup enum for game modes and rule variations
const (
	KO_SIMPLE = iota
	KO_POSITIONAL
	KO_SITUATIONAL
	KO_SPIGHT
)

// Flags for scoring modes
const (
	SCORE_AREA = iota
	SCORE_TERRITORY
)

// Flags for Taxes
const (
	TAX_NONE = iota
	TAX_SEKI
	TAX_ALL
)

// Flags for Handicap bonus for White
const (
	WHB_ZERO = iota
	WHB_N
	WHB_N_MINUS_ONE
)

const (
	KOMI_DEFAULT  = 6.5
	MIN_USER_KOMI = -150.0
	MAX_USER_KOMI = 150.0
)

type Rules struct {
	KoRule             int     `json:"ko"`                           // Ko rule to use
	ScoringRule        int     `json:"scoring"`                      // Scoring rule to use
	TaxRule            int     `json:"tax"`                          // Tax rule to use
	WhiteHandicapBonus int     `json:"whiteHandicapBonus,omitempty"` // Handicap bonus for White
	MultiStoneSuicide  bool    `json:"suicide"`                      // Allow multi-st
	HasButton          bool    `json:"hasButton"`                    // Has button
	FriendlyPassOk     bool    `json:"friendlyPassOk"`               // Friendly pass ok
	Komi               float32 `json:"komi,omitempty"`               // Komi value

}

// Constructor for declaring custom ruleset for the game
func (r *Rules) CustomRules(koRule, scoreRule, taxRule, whbRule int, suicide, button, passOk bool, komi float32) Rules {
	return Rules{
		KoRule:             koRule,
		ScoringRule:        scoreRule,
		TaxRule:            taxRule,
		WhiteHandicapBonus: whbRule,
		MultiStoneSuicide:  suicide,
		HasButton:          button,
		FriendlyPassOk:     passOk,
		Komi:               komi,
	}
}

// Construct a Tromp Taylor rule game
func (r *Rules) GetTrompTaylorish() *Rules {
	return &Rules{
		KoRule:             KO_POSITIONAL,
		ScoringRule:        SCORE_AREA,
		TaxRule:            TAX_NONE,
		MultiStoneSuicide:  false,
		HasButton:          false,
		WhiteHandicapBonus: WHB_ZERO,
		FriendlyPassOk:     false,
		Komi:               7.5,
	}
}

func (r *Rules) GetSimpleTerritory() *Rules {
	return &Rules{
		KoRule:             KO_SIMPLE,
		ScoringRule:        SCORE_TERRITORY,
		TaxRule:            TAX_SEKI,
		MultiStoneSuicide:  false,
		HasButton:          false,
		WhiteHandicapBonus: WHB_ZERO,
		FriendlyPassOk:     false,
		Komi:               7.5,
	}
}

// destructor for full interface
func (r *Rules) Close() {}

func (r *Rules) EqualsIgnoringKomi(other *Rules) bool {
	return r.KoRule == other.KoRule &&
		r.ScoringRule == other.ScoringRule &&
		r.TaxRule == other.TaxRule &&
		r.WhiteHandicapBonus == other.WhiteHandicapBonus &&
		r.MultiStoneSuicide == other.MultiStoneSuicide &&
		r.HasButton == other.HasButton &&
		r.FriendlyPassOk == other.FriendlyPassOk
}

// Checks if the final score of the game will result in an integer. This is possible
// when provided komi does not have the traditional 0.5 added to it.
func (r *Rules) GameResultWillBeInteger() bool {
	komiIsInteger := float32(int(r.Komi)) == r.Komi
	return komiIsInteger != r.HasButton
}

func komiIsIntOrHalfInt(komi float32) bool {
	return !math.IsInf(float64(komi), 0) && komi*2 == float32(int(komi)*2)
}

func koRuleStrings() []string {
	return []string{"SIMPLE", "POSITIONAL", "SITUATIONAL", "SPIGHT"}
}

func scoringRuleStrings() []string {
	return []string{"AREA", "TERRITORY"}
}

func taxRuleStrings() []string {
	return []string{"NONE", "SEKI", "ALL"}
}

func whiteHandicapBonusStrings() []string {
	return []string{"ZERO", "N", "N-1"}
}

func parseKoRule(s string) (int, error) {
	switch s {
	case "SIMPLE":
		return KO_SIMPLE, nil
	case "POSITIONAL":
		return KO_POSITIONAL, nil
	case "SITUATIONAL":
		return KO_SITUATIONAL, nil
	case "SPIGHT":
		return KO_SPIGHT, nil
	default:
		return -1, errors.New("invalid Ko Rule")
	}
}

func parseScoringRule(s string) (int, error) {
	switch s {
	case "AREA":
		return SCORE_AREA, nil
	case "TERRITORY":
		return SCORE_TERRITORY, nil
	default:
		return -1, errors.New("invalid Scoring Rule")
	}
}

func parseTaxRule(s string) (int, error) {
	switch s {
	case "NONE":
		return TAX_NONE, nil
	case "SEKI":
		return TAX_SEKI, nil
	case "ALL":
		return TAX_ALL, nil
	default:
		return -1, errors.New("invalid Tax Rule")
	}
}

func parseWhiteHandicapBonus(s string) (int, error) {
	switch s {
	case "ZERO":
		return WHB_ZERO, nil
	case "N":
		return WHB_N, nil
	case "N-1":
		return WHB_N_MINUS_ONE, nil
	default:
		return -1, errors.New("invalid White Handicap Bonus")
	}
}

func writeKoRule(koRule int) string {
	switch koRule {
	case KO_SIMPLE:
		return "SIMPLE"
	case KO_POSITIONAL:
		return "POSITIONAL"
	case KO_SITUATIONAL:
		return "SITUATIONAL"
	case KO_SPIGHT:
		return "SPIGHT"
	default:
		return "UNKNOWN"
	}
}

func writeScoringRule(scoringRule int) string {
	switch scoringRule {
	case SCORE_AREA:
		return "AREA"
	case SCORE_TERRITORY:
		return "TERRITORY"
	default:
		return "UNKNOWN"
	}
}

func writeTaxRule(taxRule int) string {
	switch taxRule {
	case TAX_NONE:
		return "NONE"
	case TAX_SEKI:
		return "SEKI"
	case TAX_ALL:
		return "ALL"
	default:
		return "UNKNOWN"
	}
}

func writeWhiteHandicapBonus(whiteHandicapBonus int) string {
	switch whiteHandicapBonus {
	case WHB_ZERO:
		return "ZERO"
	case WHB_N:
		return "N"
	case WHB_N_MINUS_ONE:
		return "N-1"
	default:
		return "UNKNOWN"
	}
}

func (r *Rules) ToStringNoKomi() string {
	var sb strings.Builder
	sb.WriteString("Ko Rule: ")
	sb.WriteString(writeKoRule(r.KoRule))
	sb.WriteString(", Scoring Rule: ")
	sb.WriteString(writeScoringRule(r.ScoringRule))
	sb.WriteString(", Tax Rule: ")
	sb.WriteString(writeTaxRule(r.TaxRule))
	sb.WriteString(", White Handicap Bonus: ")
	sb.WriteString(writeWhiteHandicapBonus(r.WhiteHandicapBonus))
	if r.MultiStoneSuicide {
		sb.WriteString(", Suicide Allowed")
	}
	if r.HasButton {
		sb.WriteString(", Has Button")
	}
	if r.WhiteHandicapBonus != WHB_ZERO {
		sb.WriteString(", White Handicap Bonus: ")
		sb.WriteString(writeWhiteHandicapBonus(r.WhiteHandicapBonus))
	}
	if r.FriendlyPassOk {
		sb.WriteString(", Friendly Pass OK")
	}
	return sb.String()
}

func (r *Rules) ToString() string {
	var sb strings.Builder
	sb.WriteString(r.ToStringNoKomi())
	sb.WriteString(", Komi: ")
	sb.WriteString(strconv.FormatFloat(float64(r.Komi), 'f', -1, 32))
	return sb.String()
}

// TODO: Is this correct? This should achieve a similar thing as https://github.com/lightvector/KataGo/blob/4dfed3ebc9dd289f52c5cb81de45bfd40af8478d/cpp/game/rules.cpp#L233 but is untested
func (r *Rules) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

func stringToBool(s string) (bool, error) {
	if s == "true" || s == "True" {
		return true, nil
	}
	if s == "false" || s == "False" {
		return false, nil
	}
	return false, errors.New("input should be 'true' or 'false'")
}

// Method to update the rules of the game. This will update in place, as well as
// return the updated rules. Note that this could fail silently if the returned
// values are not verified.
func (r *Rules) UpdateRules(k, v string) (*Rules, error) {
	switch k {
	case "ko":
		newVal, err := parseKoRule(v)
		if err != nil {
			return nil, err
		}
		r.KoRule = newVal
	case "score":
	case "scoring":
		newVal, err := parseScoringRule(v)
		if err != nil {
			return nil, err
		}
		r.ScoringRule = newVal
	case "tax":
		newVal, err := parseTaxRule(v)
		if err != nil {
			return nil, err
		}
		r.TaxRule = newVal
	case "suicide":
		newVal, err := stringToBool(v)
		if err != nil {
			return nil, err
		}
		r.MultiStoneSuicide = newVal
	case "hasButton":
		newVal, err := stringToBool(v)
		if err != nil {
			return nil, err
		}
		r.HasButton = newVal
	case "whiteHandicapBonus":
		newVal, err := parseWhiteHandicapBonus(v)
		if err != nil {
			return nil, err
		}
		r.WhiteHandicapBonus = newVal
	case "friendlyPassOk":
		newVal, err := stringToBool(v)
		if err != nil {
			return nil, err
		}
		r.FriendlyPassOk = newVal
	default:
		return nil, fmt.Errorf("%s is not a valid rule key", k)
	}
	return r, nil
}

// Original:  https://github.com/lightvector/KataGo/blob/4dfed3ebc9dd289f52c5cb81de45bfd40af8478d/cpp/game/rules.cpp#L257
// Creates a Rules object from a ruleset string. If none provided or not a valid
// rule set, will simply return TrompTaylorish Rules. Note that this differs from
// the original implementation as it does not attempt to parse a json string to
// find out the rules. WIP
// TODO: Complete this mess as time allows
func (r *Rules) parseRulesHelper(s string) *Rules {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ToLower(s)
	if s == "japanese" || s == "korean" {
		return &Rules{
			KoRule:             KO_SIMPLE,
			ScoringRule:        SCORE_TERRITORY,
			TaxRule:            TAX_SEKI,
			MultiStoneSuicide:  false,
			HasButton:          false,
			WhiteHandicapBonus: WHB_ZERO,
			FriendlyPassOk:     false,
			Komi:               6.5,
		}
	}
	if s == "chinese" {
		return &Rules{
			KoRule:             KO_SIMPLE,
			ScoringRule:        SCORE_AREA,
			TaxRule:            TAX_NONE,
			MultiStoneSuicide:  false,
			HasButton:          false,
			WhiteHandicapBonus: WHB_N,
			FriendlyPassOk:     false,
			Komi:               7.5,
		}
	}
	if s == "chineseogs" || s == "chinesekgs" {
		return &Rules{
			KoRule:             KO_POSITIONAL,
			ScoringRule:        SCORE_AREA,
			TaxRule:            TAX_NONE,
			MultiStoneSuicide:  false,
			HasButton:          false,
			WhiteHandicapBonus: WHB_N,
			FriendlyPassOk:     true,
			Komi:               7.5,
		}
	}
	if s == "ancientarea" || s == "stonescoring" {
		return &Rules{
			KoRule:             KO_SIMPLE,
			ScoringRule:        SCORE_AREA,
			TaxRule:            TAX_ALL,
			MultiStoneSuicide:  false,
			HasButton:          false,
			WhiteHandicapBonus: WHB_ZERO,
			FriendlyPassOk:     true,
			Komi:               7.5,
		}
	}
	if s == "ancientterritory" {
		return &Rules{
			KoRule:             KO_SIMPLE,
			ScoringRule:        SCORE_TERRITORY,
			TaxRule:            TAX_ALL,
			MultiStoneSuicide:  false,
			HasButton:          false,
			WhiteHandicapBonus: WHB_ZERO,
			FriendlyPassOk:     false,
			Komi:               6.5,
		}
	}
	if s == "agabutton" {
		return &Rules{
			KoRule:             KO_SITUATIONAL,
			ScoringRule:        SCORE_AREA,
			TaxRule:            TAX_NONE,
			MultiStoneSuicide:  false,
			HasButton:          true,
			WhiteHandicapBonus: WHB_N_MINUS_ONE,
			FriendlyPassOk:     true,
			Komi:               7.0,
		}
	}
	if s == "aga" || s == "bga" || s == "french" {
		return &Rules{
			KoRule:             KO_SITUATIONAL,
			ScoringRule:        SCORE_AREA,
			TaxRule:            TAX_NONE,
			MultiStoneSuicide:  false,
			HasButton:          false,
			WhiteHandicapBonus: WHB_N_MINUS_ONE,
			FriendlyPassOk:     true,
			Komi:               7.5,
		}
	}
	if s == "newzealand" || s == "nz" {
		return &Rules{
			KoRule:             KO_SITUATIONAL,
			ScoringRule:        SCORE_AREA,
			TaxRule:            TAX_NONE,
			MultiStoneSuicide:  true,
			HasButton:          false,
			WhiteHandicapBonus: WHB_ZERO,
			FriendlyPassOk:     true,
			Komi:               7.5,
		}
	}
	if s == "goe" || s == "ing" {
		return &Rules{
			KoRule:             KO_POSITIONAL,
			ScoringRule:        SCORE_AREA,
			TaxRule:            TAX_NONE,
			MultiStoneSuicide:  true,
			HasButton:          false,
			WhiteHandicapBonus: WHB_ZERO,
			FriendlyPassOk:     true,
			Komi:               7.5,
		}
	}
	return r.GetTrompTaylorish()
}

func (r *Rules) ParseRules(s string) *Rules {
	return r.parseRulesHelper(s)
}

func (r *Rules) ParseRulesWithoutKomi(s string, komi float32) *Rules {
	rules := r.parseRulesHelper(s)
	rules.Komi = komi
	return rules
}
