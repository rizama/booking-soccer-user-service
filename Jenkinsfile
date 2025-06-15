// Jenkins Pipeline untuk CI/CD User Service
// Pipeline ini akan melakukan build, test, dan deploy aplikasi secara otomatis
// Stage 1-5 -> Proses CI (Continues Integration)
// Stage 6-8 -> Proses CD (Continues Delivery/Deployment)
pipeline {
  // Menjalankan pipeline di agent manapun yang tersedia
  agent any

  // Definisi environment variables yang akan digunakan di seluruh pipeline
  environment {
    IMAGE_NAME = 'user-service'                                    // Nama Docker image
    DOCKER_CREDENTIALS = credentials('docker-credential')          // Kredensial Docker Hub dari Jenkins credentials
    GITHUB_CREDENTIALS = credentials('github-credential')          // Kredensial GitHub dari Jenkins credentials
    SSH_KEY = credentials('ssh-key')                               // SSH key untuk akses ke server remote
    HOST = credentials('host')                                     // IP/hostname server deployment
    USERNAME = credentials('username')                             // Username untuk SSH ke server
    CONSUL_HTTP_URL = credentials('consul-http-url')               // URL Consul untuk service discovery
    CONSUL_HTTP_TOKEN = credentials('consul-http-token')           // Token autentikasi Consul
    CONSUL_WATCH_INTERVAL_SECONDS = 60                             // Interval monitoring Consul dalam detik
  }

  // Definisi tahapan-tahapan dalam pipeline
  stages {
    // Stage 1: Memeriksa commit message untuk menentukan apakah pipeline harus dijalankan
    stage('Check Commit Message') {
      steps {
        script {
          // Mengambil commit message terbaru dari Git
          def commitMessage = sh(
            script: "git log -1 --pretty=%B",  // Mengambil commit message terakhir
            returnStdout: true                  // Mengembalikan output sebagai string
          ).trim()                              // Menghapus whitespace di awal dan akhir

          echo "Commit Message: ${commitMessage}"
          // Jika commit message mengandung [skip ci], maka skip pipeline
          if (commitMessage.contains("[skip ci]")) {
            echo "Skipping pipeline due to [skip ci] tag in commit message."
            currentBuild.result = 'ABORTED'     // Set status build menjadi ABORTED
            currentBuild.delete()               // Hapus build dari queue
            return                              // Keluar dari pipeline
          }

          echo "Pipeline will continue. No [skip ci] tag found in commit message."
        }
      }
    }

    // Stage 2: Menentukan target branch berdasarkan branch yang di-trigger
    stage('Set Target Branch') {
      steps {
        script {
          echo "GIT_BRANCH: ${env.GIT_BRANCH}"
          // Mapping branch origin ke target branch lokal
          if (env.GIT_BRANCH == 'origin/main') {
            env.TARGET_BRANCH = 'main'          // Set target branch ke main
          } else if (env.GIT_BRANCH == 'origin/development') {
            env.TARGET_BRANCH = 'development'   // Set target branch ke development
          }

          echo "TARGET_BRANCH: ${env.TARGET_BRANCH}"
        }
      }
    }

    // Stage 3: Checkout kode dari repository GitHub (Mengambil/Pull source code ke Jenkins untuk di-build, test, dan deploy)
    stage('Checkout Code') {
      steps {
        script {
          // URL repository GitHub
          def repoUrl = 'https://github.com/rizama/booking-soccer-user-service.git'

          // Checkout kode dari branch yang ditentukan
          checkout([$class: 'GitSCM',
            branches: [
              [name: "*/${env.TARGET_BRANCH}"]    // Checkout branch sesuai TARGET_BRANCH
            ],
            userRemoteConfigs: [
              [url: repoUrl, credentialsId: 'github-credential']  // Menggunakan kredensial GitHub
            ]
          ])

          // Menampilkan isi direktori setelah checkout
          sh 'ls -lah'
        }
      }
    }

    // Stage 4: Login ke Docker Hub untuk push image
    stage('Login to Docker Hub') {
      steps {
        script {
          // Menggunakan kredensial Docker Hub yang tersimpan di Jenkins
          // Login ke Docker Hub menggunakan password dari stdin untuk keamanan
          withCredentials([usernamePassword(credentialsId: 'docker-credential', passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
            sh """
            echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
            """
          }
        }
      }
    }

    // Stage 5: Build Docker image dan push ke Docker Hub
    stage('Build and Push Docker Image') {
      steps {
        script {
          // Menggunakan build number Jenkins sebagai tag image
          def runNumber = currentBuild.number
          // Build Docker image dengan tag berdasarkan build number
          sh "docker build -t ${DOCKER_CREDENTIALS_USR}/${IMAGE_NAME}:${runNumber} ."
          // Push image ke Docker Hub
          sh "docker push ${DOCKER_CREDENTIALS_USR}/${IMAGE_NAME}:${runNumber}"
        }
      }
    }

    // Stage 6: Update file docker-compose.yaml dengan image tag terbaru
    stage('Update docker-compose.yaml') {
      steps {
        script {
          def runNumber = currentBuild.number
          // Menggunakan sed untuk mengganti tag image di docker-compose.yaml
          // Pattern regex mencari image dengan tag angka dan menggantinya dengan build number terbaru
          sh "sed -i 's|image: ${DOCKER_CREDENTIALS_USR}/${IMAGE_NAME}:[0-9]\\+|image: ${DOCKER_CREDENTIALS_USR}/${IMAGE_NAME}:${runNumber}|' docker-compose.yml"
        }
      }
    }

    // Stage 7: Commit dan push perubahan docker-compose.yaml ke repository
    stage('Commit and Push Changes') {
      steps {
          // Konfigurasi Git user untuk commit
          // Set remote URL dengan kredensial untuk push
          // Add file yang diubah ke staging area
          // Commit dengan message yang mengandung [skip ci] untuk menghindari trigger pipeline lagi
          // Pull latest changes dengan rebase untuk menghindari merge conflict
          // Push perubahan ke repository
        script {
          sh """
          git config --global user.name 'Jenkins CI'
          git config --global user.email 'jenkins@company.com'
          git remote set-url origin https://${GITHUB_CREDENTIALS_USR}:${GITHUB_CREDENTIALS_PSW}@github.com/rizama/booking-soccer-user-service.git
          git add docker-compose.yml
          git commit -m 'Update image version to ${TARGET_BRANCH}-${currentBuild.number} [skip ci]' || echo 'No changes to commit'
          git pull origin ${TARGET_BRANCH} --rebase
          git push origin HEAD:${TARGET_BRANCH}
          """
        }
      }
    }

    // Stage 8: Deploy aplikasi ke server remote
    stage('Deploy to Remote Server') {
      steps {
        script {
          // Direktori target di server remote
          def targetDir = "/home/rizkysamp/booking-soccer/user-service"
          // Command SSH untuk deployment
          // SSH ke server dengan key authentication, disable host key checking
          // Cek apakah direktori sudah ada dan merupakan Git repository
            // Pull latest changes dari branch yang ditentukan
            // Clone repository jika direktori belum ada

          // Copy file .env.example menjadi .env untuk konfigurasi
          // Update konfigurasi environment variables di file .env

          // Jalankan Docker Compose untuk deploy aplikasi
          // --build: rebuild image, --force-recreate: recreate container, -d: detached mode

          def sshCommandToServer = """
          ssh -o StrictHostKeyChecking=no -i ${SSH_KEY} ${USERNAME}@${HOST} '
            if [ -d "${targetDir}/.git" ]; then
                echo "Directory exists. Pulling latest changes."
                cd "${targetDir}"
                git pull origin "${TARGET_BRANCH}"
            else
                echo "Directory does not exist. Cloning repository."
                git clone -b "${TARGET_BRANCH}" git@github.com:rizama/booking-soccer-user-service.git "${targetDir}"
                cd "${targetDir}"
            fi

            cp .env.example .env
            sed -i "s/^TIMEZONE=.*/TIMEZONE=Asia\\/Jakarta/" "${targetDir}/.env"
            sed -i "s/^CONSUL_HTTP_URL=.*/CONSUL_HTTP_URL=${CONSUL_HTTP_URL}/" "${targetDir}/.env"
            sed -i "s/^CONSUL_HTTP_PATH=.*/CONSUL_HTTP_PATH=backend\\/user-service/" "${targetDir}/.env"
            sed -i "s/^CONSUL_HTTP_TOKEN=.*/CONSUL_HTTP_TOKEN=${CONSUL_HTTP_TOKEN}/" "${targetDir}/.env"
            sed -i "s/^CONSUL_WATCH_INTERVAL_SECONDS=.*/CONSUL_WATCH_INTERVAL_SECONDS=${CONSUL_WATCH_INTERVAL_SECONDS}/" "${targetDir}/.env"

            sudo docker compose up -d --build --force-recreate
          '
          """
          // Eksekusi command SSH
          sh sshCommandToServer
        }
      }
    }
  }
}
