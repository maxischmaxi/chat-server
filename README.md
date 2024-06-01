# Deployment

- gcloud compute instances create chat-server --zone=europe-west10-a --machine-type=e2-medium --image-family=debian-10 --image-project=debian-cloud
- gcloud compute firewall-rules create allow-tcp-rule --allow=tcp:8080 --source-ranges=0.0.0.0/0 --target-tags=tcp-server
- gcloud compute instances add-tags chat-server --tags=tcp-server
- gcloud compute ssh chat-server # start the server with git clone etc
- gcloud compute instances list

