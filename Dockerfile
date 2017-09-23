FROM ruby
WORKDIR /

# Prepare
RUN apt-get update
RUN apt-get dist-upgrade -y
RUN apt-get install -y -u apt-utils unzip nodejs upstart wget curl jruby cron nano xvfb screen htop npm openssl git libgtk2.0-0 libxtst6 libxss1 libgconf-2-4 libnss3
RUN gem install aquatone
ENTRYPOINT ["/bin/bash"]
