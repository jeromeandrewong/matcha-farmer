provider "aws" {
  region = "ap-southeast-1"
  profile = "default"
}

variable "lambda_role_arn" {
  type = string
}

variable "telegram_bot_token" {
  type        = string
  description = "Telegram Bot API Token"
}

variable "telegram_chat_id" {
  type        = string
  description = "Telegram Chat ID to send messages to"
}

resource "aws_lambda_function" "matcha-farmer" {
  function_name = "matcha-farmer"
  filename      = "../lambda.zip"
  source_code_hash = filebase64sha256("../lambda.zip")
  handler = "main"
  runtime = "provided.al2"
  role = var.lambda_role_arn
  timeout = 30
  memory_size = 128

  environment {
    variables = {
      TELEGRAM_BOT_TOKEN = var.telegram_bot_token
      TELEGRAM_CHAT_ID   = var.telegram_chat_id
    }
  }
}
