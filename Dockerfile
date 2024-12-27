FROM debian:stable-slim

COPY color_my_practice /bin/color_my_practice

COPY swagger/ swagger/  

CMD [ "/bin/color_my_practice" ]