NAME=MARKBOT
DOCKER=docker compose

clean:
	rm -rf api/tmp/ && rm api/src/database/toggles.json

api-console:
	$(DOCKER) run --service-ports api-markbot bash
remove-containers:
	$(DOCKER) down --remove-orphans

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
