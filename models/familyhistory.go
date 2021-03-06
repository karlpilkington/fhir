// Copyright (c) 2011-2014, HL7, Inc & The MITRE Corporation
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
//     * Redistributions of source code must retain the above copyright notice, this
//       list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright notice,
//       this list of conditions and the following disclaimer in the documentation
//       and/or other materials provided with the distribution.
//     * Neither the name of HL7 nor the names of its contributors may be used to
//       endorse or promote products derived from this software without specific
//       prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package models

import "time"

type FamilyHistory struct {
	Id         string                           `json:"-" bson:"_id"`
	Identifier []Identifier                     `bson:"identifier,omitempty", json:"identifier,omitempty"`
	Subject    Reference                        `bson:"subject,omitempty", json:"subject,omitempty"`
	Date       FHIRDateTime                     `bson:"date,omitempty", json:"date,omitempty"`
	Note       string                           `bson:"note,omitempty", json:"note,omitempty"`
	Relation   []FamilyHistoryRelationComponent `bson:"relation,omitempty", json:"relation,omitempty"`
}

// This is an ugly hack to deal with embedded structures in the spec condition
type FamilyHistoryRelationConditionComponent struct {
	Type        CodeableConcept `bson:"type,omitempty", json:"type,omitempty"`
	Outcome     CodeableConcept `bson:"outcome,omitempty", json:"outcome,omitempty"`
	OnsetAge    Quantity        `bson:"onsetAge,omitempty", json:"onsetAge,omitempty"`
	OnsetRange  Range           `bson:"onsetRange,omitempty", json:"onsetRange,omitempty"`
	OnsetString string          `bson:"onsetString,omitempty", json:"onsetString,omitempty"`
	Note        string          `bson:"note,omitempty", json:"note,omitempty"`
}

// This is an ugly hack to deal with embedded structures in the spec relation
type FamilyHistoryRelationComponent struct {
	Name            string                                    `bson:"name,omitempty", json:"name,omitempty"`
	Relationship    CodeableConcept                           `bson:"relationship,omitempty", json:"relationship,omitempty"`
	BornPeriod      Period                                    `bson:"bornPeriod,omitempty", json:"bornPeriod,omitempty"`
	BornDate        FHIRDateTime                              `bson:"bornDate,omitempty", json:"bornDate,omitempty"`
	BornString      string                                    `bson:"bornString,omitempty", json:"bornString,omitempty"`
	AgeAge          Quantity                                  `bson:"ageAge,omitempty", json:"ageAge,omitempty"`
	AgeRange        Range                                     `bson:"ageRange,omitempty", json:"ageRange,omitempty"`
	AgeString       string                                    `bson:"ageString,omitempty", json:"ageString,omitempty"`
	DeceasedBoolean *bool                                     `bson:"deceasedBoolean,omitempty", json:"deceasedBoolean,omitempty"`
	DeceasedAge     Quantity                                  `bson:"deceasedAge,omitempty", json:"deceasedAge,omitempty"`
	DeceasedRange   Range                                     `bson:"deceasedRange,omitempty", json:"deceasedRange,omitempty"`
	DeceasedDate    FHIRDateTime                              `bson:"deceasedDate,omitempty", json:"deceasedDate,omitempty"`
	DeceasedString  string                                    `bson:"deceasedString,omitempty", json:"deceasedString,omitempty"`
	Note            string                                    `bson:"note,omitempty", json:"note,omitempty"`
	Condition       []FamilyHistoryRelationConditionComponent `bson:"condition,omitempty", json:"condition,omitempty"`
}

type FamilyHistoryBundle struct {
	Type         string                     `json:"resourceType,omitempty"`
	Title        string                     `json:"title,omitempty"`
	Id           string                     `json:"id,omitempty"`
	Updated      time.Time                  `json:"updated,omitempty"`
	TotalResults int                        `json:"totalResults,omitempty"`
	Entry        []FamilyHistoryBundleEntry `json:"entry,omitempty"`
	Category     FamilyHistoryCategory      `json:"category,omitempty"`
}

type FamilyHistoryBundleEntry struct {
	Title    string                `json:"title,omitempty"`
	Id       string                `json:"id,omitempty"`
	Content  FamilyHistory         `json:"content,omitempty"`
	Category FamilyHistoryCategory `json:"category,omitempty"`
}

type FamilyHistoryCategory struct {
	Term   string `json:"term,omitempty"`
	Label  string `json:"label,omitempty"`
	Scheme string `json:"scheme,omitempty"`
}
