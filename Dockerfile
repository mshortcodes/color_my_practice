FROM debian:stable-slim

ADD color_my_practice /bin/color_my_practice

CMD [ "/bin/color_my_practice" ]