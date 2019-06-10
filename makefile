db:
	ssh -nNT -L "0.0.0.0:3309:koorong-master.c2k4u1dberoy.ap-northeast-2.rds.amazonaws.com:3306" ec2-user@13.209.6.170 -i ~/.ssh/TylerKeyPair.pem

docker:
	docker build -t tylergrey/lotte-server .
	docker push tylergrey/lotte-server:latest

.PHONY: db docker