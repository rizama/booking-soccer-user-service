# Definisi variabel warna untuk output terminal yang lebih menarik
GREEN  := $(shell tput -Txterm setaf 2)  # Warna hijau untuk pesan sukses
YELLOW := $(shell tput -Txterm setaf 3)  # Warna kuning untuk peringatan
WHITE  := $(shell tput -Txterm setaf 7)  # Warna putih untuk teks normal
CYAN   := $(shell tput -Txterm setaf 6)  # Warna cyan untuk informasi
RESET  := $(shell tput -Txterm sgr0)     # Reset warna ke default

## Live reload:
# Target untuk mempersiapkan tools yang diperlukan untuk hot reload
watch-prepare: ## Install the tools required for the watch command
# Download dan install Air (hot reload tool untuk Go) dari GitHub
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh

# Target untuk menjalankan service dengan hot reload menggunakan Air
watch: ## Run the service with hot reload
# Menjalankan Air dari direktori bin/ untuk hot reload development
	bin/air

## Build:
# Target untuk build aplikasi Go menjadi binary executable
build: ## Build the service
# Compile Go code menjadi binary dengan nama 'user-service'
	go build -o user-service

## Docker:
# Target untuk menjalankan service menggunakan Docker Compose
# Menjalankan docker-compose dengan rebuild image dan recreate container
# -d: detached mode (background), --build: rebuild image, --force-recreate: recreate container
docker-compose: ## Start the service in docker
	docker-compose up -d --build --force-recreate

# Target untuk build Docker image dengan tag yang ditentukan
docker-build: ## Build the Docker image with a specified tag
# Menampilkan pesan informasi dengan warna cyan
	@echo "$(CYAN)Building Docker image...$(RESET)"
# Validasi apakah parameter 'tag' sudah diberikan
# Jika tag kosong, tampilkan error dengan warna kuning dan keluar
	@if [ -z "$(tag)" ]; then \
		echo "$(YELLOW)Error: Please specify the 'tag' parameter, e.g., make docker-build tag=1.0.0$(RESET)"; \
		exit 1; \
	fi
# Build Docker image dengan platform linux/amd64 dan tag yang ditentukan
	docker build --platform linux/amd64 -t rizkysamp/bwa-soccer-user-service:$(tag) .
# Menampilkan pesan sukses dengan warna hijau
	@echo "$(GREEN)Docker image built with tag '$(tag)'$(RESET)"

# Target untuk push Docker image ke registry (Docker Hub)
docker-push: ## Push the Docker image with a specified tag
# Menampilkan pesan informasi dengan warna cyan
	@echo "$(CYAN)Pushing Docker image...$(RESET)"
# Validasi apakah parameter 'tag' sudah diberikan
# Jika tag kosong, tampilkan error dengan warna kuning dan keluar
	@if [ -z "$(tag)" ]; then \
		echo "$(YELLOW)Error: Please specify the 'tag' parameter, e.g., make docker-push tag=1.0.0$(RESET)"; \
		exit 1; \
	fi
# Push Docker image ke registry dengan tag yang ditentukan
	docker push rizkysamp/bwa-soccer-user-service:$(tag)
# Menampilkan pesan sukses dengan warna hijau
	@echo "$(GREEN)Docker image pushed with tag '$(tag)'$(RESET)"
