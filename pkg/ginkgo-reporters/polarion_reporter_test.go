/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Red Hat, Inc.
 *
 */

package ginkgo_reporters

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/config"
	"github.com/onsi/ginkgo/v2/types"
	. "github.com/onsi/gomega"
)

var _ = Describe("ginkgo_reporters", func() {

	Context("When Spec Suite Will Begin", func() {

		var (
			suiteSummary *types.SuiteSummary
			configType   config.GinkgoConfigType
			properties   PolarionProperties
			reporter     PolarionReporter
		)
		suiteSummary = &types.SuiteSummary{
			SuiteDescription: "SUITE DESCRIPTION",
		}

		configType = config.GinkgoConfigType{
			ParallelTotal: 1,
			ParallelNode:  1,
		}

		reporter = PolarionReporter{
			Run:             true,
			Filename:        "polarion.xml",
			ProjectId:       "QE",
			PlannedIn:       "QE_1.0",
			LookupMethod:    "id",
			TestSuiteParams: "OS=EL7 SC=NFS,ENV_TIER=TIER1",
			TestRunTemplate: "QE_Regression_template",
			TestRunTitle:    "My run title",
		}

		properties = PolarionProperties{
			Property: []PolarionProperty{
				{
					Name:  "polarion-custom-OS",
					Value: "EL7",
				},
				{
					Name:  "polarion-custom-SC",
					Value: "NFS",
				},
				{
					Name:  "polarion-custom-ENV_TIER",
					Value: "TIER1",
				},
				{
					Name:  "polarion-project-id",
					Value: "QE",
				},
				{
					Name:  "polarion-lookup-method",
					Value: "id",
				},
				{
					Name:  "polarion-custom-plannedin",
					Value: "QE_1.0",
				},
				{
					Name:  "polarion-testrun-id",
					Value: "QE_1.0_EL7_NFS_TIER1",
				},
				{
					Name:  "polarion-custom-isautomated",
					Value: "True",
				},
				{
					Name:  "polarion-testrun-status-id",
					Value: "inprogress",
				},
				{
					Name:  "polarion-testrun-template-id",
					Value: "QE_Regression_template",
				},
				{
					Name:  "polarion-testrun-title",
					Value: "My run title",
				},
			},
		}

		It("Should info reporter test suite name & properties", func() {
			reporter.SuiteWillBegin(configType, suiteSummary)
			Expect(reporter.TestSuiteName).To(Equal(suiteSummary.SuiteDescription))
			Expect(reporter.Suite.Properties).To(Equal(properties))
		})
	})

	Context("When Spec Did Complete", func() {

		var (
			specSummaryPass *types.SpecSummary
			specSummaryFail *types.SpecSummary
			specSummarySkip *types.SpecSummary
			reporter        PolarionReporter
		)
		specSummaryPass = &types.SpecSummary{
			ComponentTexts: []string{"THIS IS", "[test_id:123]A PASSING", "TEST"},
			State:          types.SpecStatePassed,
			CapturedOutput: "Test output",
		}

		specSummaryFail = &types.SpecSummary{
			ComponentTexts: []string{"THIS IS", "A FAILING", "TEST"},
			State:          types.SpecStateFailed,
			CapturedOutput: "Test output",
			Failure: types.SpecFailure{
				Message: "ERROR MSG",
				Location: types.CodeLocation{
					FileName:       "file/a",
					LineNumber:     3,
					FullStackTrace: "some-stack-trace",
				},
				ComponentCodeLocation: types.CodeLocation{
					FileName:       "file/b",
					LineNumber:     4,
					FullStackTrace: "some-stack-trace",
				},
			},
		}

		skipMessage := JUnitSkipped{}

		specSummarySkip = &types.SpecSummary{
			ComponentTexts: []string{"THIS IS", "[test_id:789]A SKIPPING", "TEST"},
			State:          types.SpecStateSkipped,
			CapturedOutput: "Test output",
		}

		testcases := []PolarionTestCase{
			{
				Name: fmt.Sprintf("%s: %s", "THIS IS", "[test_id:123]A PASSING TEST"),
				Properties: PolarionProperties{
					Property: []PolarionProperty{
						{
							Name:  "polarion-testcase-id",
							Value: "PREFIX-123",
						},
					},
				},
			},
			{
				Name:       fmt.Sprintf("%s: %s", "THIS IS", "A FAILING TEST"),
				Properties: PolarionProperties{},
				SystemOut:  "Test output",
				FailureMessage: &JUnitFailureMessage{
					Type:    "Failure",
					Message: "file/b:4\nERROR MSG\nfile/a:3",
				},
			},
			{
				Name: fmt.Sprintf("%s: %s", "THIS IS", "[test_id:789]A SKIPPING TEST"),
				Properties: PolarionProperties{
					Property: []PolarionProperty{
						{
							Name:  "polarion-testcase-id",
							Value: "PREFIX-789",
						},
					},
				},
				Skipped: &skipMessage,
			},
		}

		reporter = PolarionReporter{
			Run:          true,
			Filename:     "polarion.xml",
			ProjectId:    "QE",
			PlannedIn:    "QE_1.0",
			TestIDPrefix: "PREFIX",
		}

		It("Should info reporter test cases did complete", func() {
			reporter.SpecDidComplete(specSummaryPass)
			reporter.SpecDidComplete(specSummaryFail)
			reporter.SpecDidComplete(specSummarySkip)
			Expect(reporter.Suite.TestCases).To(Equal(testcases))
		})
	})

	Context("When Spec Spec Suite Did End", func() {

		var (
			suiteSummary *types.SuiteSummary
			reporter     PolarionReporter
		)
		suiteSummary = &types.SuiteSummary{
			NumberOfSpecsThatWillBeRun: 3,
			RunTime:                    time.Minute,
			NumberOfFailedSpecs:        1,
		}

		reporter = PolarionReporter{
			Run:             true,
			Filename:        "polarion.xml",
			ProjectId:       "QE",
			PlannedIn:       "QE_1.0",
			TestSuiteParams: "",
		}

		It("Should info number of tests, execution time, number of failures & verify report file was created", func() {
			reporter.SuiteDidEnd(suiteSummary)
			Expect(reporter.Suite.Tests).To(Equal(3))
			Expect(reporter.Suite.Time).To(Equal(time.Minute.Seconds()))
			Expect(reporter.Suite.Failures).To(Equal(1))
			if _, err := os.Stat(reporter.Filename); err == nil {
				Expect(err).ToNot(HaveOccurred())
			}
		})
	})

})
