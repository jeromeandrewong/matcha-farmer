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

resource "aws_cloudwatch_event_rule" "matcha_farmer_schedule" {
  name                = "matcha-farmer-schedule"
  description         = "Schedule for Matcha Farmer Lambda Function"
  schedule_expression = "cron(0 4,16 * * ? *)"# every 12 hours at midnight and noon SGT
}

resource "aws_cloudwatch_event_target" "matcha_farmer_target" {
  rule      = aws_cloudwatch_event_rule.matcha_farmer_schedule.name
  target_id = "matcha-farmer"
  arn       = aws_lambda_function.matcha-farmer.arn
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.matcha-farmer.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.matcha_farmer_schedule.arn
}
