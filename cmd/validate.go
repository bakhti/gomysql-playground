/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	_ "github.com/bakhti/gomysql-playground/pkg/validators"
	_ "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate data between databases and tables",
	Long:  `This is an experiment where 2 databases will be compared.`,
	Run: func(cmd *cobra.Command, args []string) {
		// log := logrus.New()
		// log.SetLevel(logrus.DebugLevel)
		// logger := logrus.NewEntry(log)
		// validator := validators.NewValidator(logger)
		// if err := validator.Run(); err != nil {
		// 	logger.WithError(err).Fatal("Could not validate data")
		// }
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
