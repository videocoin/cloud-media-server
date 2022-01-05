package mediacore

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func MP4BoxDashExec(inputVideoPath, inputAudioPath, outputPath string) (string, error) {
	ctx := context.Background()

	cmdArgs := []string{
		"-dash", "-1", "-bs-switching", "no", "-single-file",
		"-segment-name", "%s", "-url-template",
		"-out", outputPath,
	}

	if inputAudioPath != "" {
		cmdArgs = append(cmdArgs, "-rap")
		cmdArgs = append(cmdArgs, inputAudioPath)
	}

	cmdArgs = append(cmdArgs, inputVideoPath)

	fmt.Println("MP4Box " + strings.Join(cmdArgs, " "))

	cmd := exec.CommandContext(ctx, "MP4Box", cmdArgs...)
	out, err := cmd.CombinedOutput()
	outStr := string(out)
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), outStr)
	}

	return outStr, nil
}

func MP4BoxCryptExec(drmXmlPath, inputPath, outputPath string) (string, error) {
	ctx := context.Background()
	cmdArgs := []string{"-crypt", drmXmlPath, inputPath, "-out", outputPath}

	fmt.Println("MP4Box " + strings.Join(cmdArgs, " "))

	cmd := exec.CommandContext(ctx, "MP4Box", cmdArgs...)
	out, err := cmd.CombinedOutput()
	outStr := string(out)
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), outStr)
	}

	return outStr, nil
}
