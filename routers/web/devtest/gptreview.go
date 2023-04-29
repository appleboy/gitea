package devtest

import (
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/gpt"
	"errors"
	"net/http"
	"strings"
)

func gptSimpleReview(ctx *context.Context, content string) (string, error) {
	c, err := gpt.NewOpenAIClientFromSetting()
	if err != nil {
		return "", err
	}

	resp, err := gpt.CreateSimpleChatCompletion(ctx, c, content)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func GptReview(ctx *context.Context) {
	reviewType := ctx.FormString("type")

	var respContent string
	var respErr error
	if reviewType == "diff" {
		diff := ctx.FormString("diff")

		content := "You are an experienced developer, and you are reviewing the following git patch. Please provide review suggestions for it.\n\n"
		content += diff

		respContent, respErr = gptSimpleReview(ctx, diff)
	} else if reviewType == "comment" {
		comments := ctx.Req.Form["comments"]

		content := "You are an experienced developer, and you are reading the comments for a git patch. Please summarize them.\n\n"
		content += strings.Join(comments, "\n")
		respContent, respErr = gptSimpleReview(ctx, content)
	} else if reviewType == "code" {
		code := ctx.FormString("code")

		content := "You are an experienced developer, please complement the code for the comment TODO, and only output the code. The original code is:\n\n"
		content += code

		respContent, respErr = gptSimpleReview(ctx, content)
	} else {
		ctx.ServerError("UnknownReviewType", errors.New("unknown review type"))
		return
	}

	if respErr != nil {
		ctx.ServerError("GptReviewError", respErr)
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"reviewResponse": respContent,
	})
}
