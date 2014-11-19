Instalación de postfix con relay a mandrill 
=====

###Crear cuenta en mandrill app 

1) Dirigirse a la sección “settings” -> Domains 
2) Dentro de esta sección agregar el un dominio de salida
3) Dentro de esta sección dar clic al botón “View DKIM/SPF setup instructions”
4) Crear los registros de tipo TXT en el administrador de DNS del proveedor de servicios 
	
Registro SPF 
	```	
		Crear registro TXT en blanco con valor especificado por mandrill 
	```

Registro DKIM
	```	 
		Crear registro TXT con este nombre mandrill._domainkey con valor especificado por Mandril  
	```

Dar clic en verificar DNS

###Preparar servidor Linux (CentOS 6 64bits)

	#####Actualizar

	```
		yum update
	```

	#####Desinstalar sendmail 

	```
		yum remove sendmail
	```

	#####Instalar postfix y cyrus

	```
		yum install postfix
	```

	```
		yum install cyrus-sasl-plain
	```

	#####Definir postfix como auto arranque

	```
		chkconfig postfix on
	```

	#####Crear fichero con credenciales /etc/postfix/sasl_passwd con siguiente contenido 

	```
		[smtp.mandrillapp.com] NombreUsuarioMandrill:API_KEY
	```

	```
		chmod 600 /etc/postfix/sasl_passwd
	```

	#####Mapear fichero con credenciales

	```
		postmap /etc/postfix/sasl_passwd
	```	

	#####Agregar siguientes líneas al fichero /etc/postfix/main.cf

	```
		smtp_sasl_auth_enable = yes
		smtp_sasl_password_maps = hash:/etc/postfix/sasl_passwd 
		smtp_sasl_security_options = noanonymous
		smtp_use_tls = yes 
		relayhost = [smtp.mandrillapp.com]
	```

	#####Iniciar postfix

	```
		service postfix start
	```

	#####Enviar correo de prueba

	```
		$ sendmail destino@dominio.com
		From: nombre@tudominio.com
		Subject: xxxxxxxx
		cuerpo del mensaje
		.
	```
