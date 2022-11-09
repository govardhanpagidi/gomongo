FROM golang:latest

LABEL version="1.0.0"

LABEL "maintainer"="Jason Mimick <jason.mimick@mongodb.com>"
LABEL "repository"="https://github.com/mongodb/mongodbatlas-cloudformation-resources"
LABEL "original-credit"="https://github.com/mridhul/cfn-deploy"

LABEL com.github.actions.name="MongoDB Atlas Cloudformation Github Deploy"
LABEL com.github.actions.description="Deploy the MongoDB Atlas Cloudformation custom resources to any AWS region with a Github action."
LABEL com.github.actions.icon="upload-cloud"
LABEL com.github.actions.color="green"

RUN apt-get update && apt-get install -y awscli git build-essential python3-pip 
RUN rm -rf /var/lib/apt/lists/*
RUN pip3 install cloudformation-cli cloudformation-cli-java-plugin cloudformation-cli-go-plugin cloudformation-cli-python-plugin
COPY requirements.txt /
RUN pip3 install -r /requirements.txt
RUN apt-get install jq

# Install mongocli
RUN MCLI_TAG=$(curl -sL --header "Accept: application/json" https://github.com/mongodb/mongocli/releases/latest | jq -r '.["tag_name"]') && \
    MCLI_VERSION=$(echo $MCLI_TAG | cut -dv -f2) && \
    MCLI_DEB="mongocli_${MCLI_VERSION}_linux_x86_64.deb" && \
    curl -OL https://github.com/mongodb/mongocli/releases/download/${MCLI_TAG}/${MCLI_DEB} && \
    echo "About to install mongocli from: ${MCLI_DEB}" && \
    dpkg -i ${MCLI_DEB} \


#WORKDIR /atlas-cfn
#COPY atlas-cfn-deploy.py atlas-cfn-deploy.py  
#COPY atlas-cfn-stack-cleaner.sh atlas-cfn-stack-cleaner.sh
#RUN chmod +x atlas-cfn-stack-cleaner.sh

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
