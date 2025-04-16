NAME=MARKBOT
DOCKER=docker-compose

# Очистка временных файлов
clean:
	rm -rf api/tmp/ && rm -f api/src/database/toggles.json

# Доступ к консоли в dev-контейнере
api-console-dev:
	$(DOCKER) --profile dev run --service-ports api-markbot-dev bash

# Доступ к консоли в prod-контейнере
api-console-prod:
	$(DOCKER) --profile prod run --service-ports api-markbot-prod bash

# Удаление всех контейнеров
remove-containers:
	$(DOCKER) down --remove-orphans

# Запуск в режиме разработки с hot-reloading (air)
dev:
	$(DOCKER) --profile dev up

# Запуск в режиме разработки с отсоединением (detached mode)
dev-detached:
	$(DOCKER) --profile dev up -d

# Сборка для режима разработки
build-dev:
	$(DOCKER) --profile dev build

# Запуск в production режиме
prod:
	$(DOCKER) --profile prod up

# Запуск в production режиме с отсоединением (detached mode)
prod-detached:
	$(DOCKER) --profile prod up -d

# Сборка для production режима
build-prod:
	$(DOCKER) --profile prod build

# Запуск всех сервисов (dev и prod)
all:
	$(DOCKER) --profile all up

# Приветственное сообщение
greeting:
	@echo "🎓🎓🎓"
	@echo
	@echo "Hello, dear! My name is $(NAME). I'm a bot and web-app (and your best companion)🦾👷‍♀️"
	@echo "You can explore my guts for learning a lot of useful information^^"
	@echo "To some, I may seem like a complicated girl. But really, I'm as simple as a daisy 🌼"
	@echo "Start studying me, let's be friends!❤️‍🔥"
	@echo
	@echo "Go to the README.md file or Wiki on Github for more information."
	@echo
	@echo "🎓🎓🎓"
