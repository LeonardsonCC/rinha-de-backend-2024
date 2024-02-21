# RINHA

## v1.0 - 100% puro (exceto pelo driver postgres)
Os testes de carga não estavam passando, mas acredito que era por eu não ter recriado a base de teste (pois é, não sou uma pessoa inteligente).

## v1.1 - Adicionado [pgx](https://github.com/jackc/pgx/)
Adicionado o uso de [pgx](https://github.com/jackc/pgx/) para não precisar fazer minha pool de conexões.

Teste de carga:
![image](https://github.com/LeonardsonCC/rinha-de-backend-2024/assets/21212048/d9b7e2f3-74eb-4e1c-8dcb-6789edbfe86b)
