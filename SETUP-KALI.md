# Kali Linux VM Setup

Run these commands inside the Kali VM, not in Windows PowerShell.

## Base utilities

```bash
sudo apt update
sudo apt install -y git curl wget unzip jq make ca-certificates gnupg
```

## Docker

```bash
sudo apt install -y docker.io
sudo systemctl enable docker --now
sudo usermod -aG docker "$USER"
```

Log out of Kali and log back in, then verify:

```bash
docker --version
docker run --rm hello-world
```

## kubectl

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256"
echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
rm kubectl kubectl.sha256
kubectl version --client
```

## kind

```bash
curl -Lo kind https://kind.sigs.k8s.io/dl/v0.32.0/kind-linux-amd64
chmod +x kind
sudo mv kind /usr/local/bin/kind
kind version
```

## Terraform

```bash
TERRAFORM_VERSION=1.15.7
curl -LO "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
unzip "terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
sudo install -o root -g root -m 0755 terraform /usr/local/bin/terraform
rm terraform "terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
terraform version
```

## Helm

```bash
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-4
chmod 700 get_helm.sh
./get_helm.sh
rm get_helm.sh
helm version
```

## Final verification

```bash
git --version
docker --version
kubectl version --client
kind version
terraform version
helm version
```
