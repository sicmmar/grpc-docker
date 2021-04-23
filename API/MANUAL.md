# Api con NodeJs

Contenido
- [Api Rest](#api-rest)
    - [Endpoints](#endpoints) 
- [Dockerizando](#dockerizando)
    - [Dockerfile](#dockerfile)
    - [Docker Compose](#docker-compose)

### Api Rest
Las dependencias para levantar la Api Rest con Node, son:
```json
"dependencies": {
    "body-parser": "^1.19.0",
    "cors": "^2.8.5",
    "express": "^4.17.1",
    "express-force-https": "^1.0.0",
    "mongodb": "^3.6.4"
}
```

Primero se conecta a la base de datos, y se defina un ```GET``` raíz que tiene un mensaje de bienvenida.

```js
mongoClient.connect(urlMongo, { useUnifiedTopology: true })
.then(client => {
    console.log("Conectado a la base de datos!")
    const db = client.db(nameDB)
    const coleccion = db.collection('usuario')

    app.get('/', (req, res) => {
        console.log('inicio de api')
        res.send('API SOPES 1 :D');
    });

    //-----------------------------------------------------------------------------------
    // ------------------ ACÁ SE COLOCAN TODOS LOS ENDPOINTS ----------------------------
    //-----------------------------------------------------------------------------------

    app.listen(port, () => {console.log(`Server corriendo en puerto ${port}!`) });
    
})
.catch(console.error)
```
A continuación se definen los endpoints utilizados.

#### Endpoints
- Módulo de Procesos

```js
app.get('/procesos', (req, res) => {
    const html = fs.readFileSync('/elements/procs/procesos','utf-8');
    let texto = html.toString();
    console.log("Procesos!")
    res.send(texto);
});
```
- Módulo RAM
```js
app.get('/ram', (req, res) => {
    const html = fs.readFileSync('/elements/procs/ram-module','utf-8');
    let texto = html.toString();
    console.log("Ram!")
    res.send(texto);
});
```
- Nuevo registro en Mongo

Para este endpoint, se verifica que ningún registro sea nulo para que se pueda registrar en la base de datos.
```js
app.post('/nuevoRegistro', (req, res) => {
    const data = req.body;
    
    if(data.name == null || data.location == null || data.age == null || data.vaccine_type == null || data.gender == null ||
        data.name == "" || data.location == "" || data.age == 0 || data.vaccine_type == "" || data.gender == "")
    {
        res.send('Nulls encontrados');
    }
    else{
        const user = {
            "name": data.name,
            "location": data.location,
            "gender": data.gender,
            "age": data.age,
            "vaccine_type": data.vaccine_type,
            "way": data.way
        }

        coleccion.insertOne(user)
        .then(result => {
            console.log("Registro Insertado!");
            res.send('Registro Insertado!');
        })
        .catch(error => console.error("Error al insertar un registro: ", error));
    }
});
```
- Eliminar todos los registros en Mongo

Con la opción drop de Mongo, se elimnan los registros existentes.
```js
app.get('/deleteAll', (req, res) => {
    coleccion.drop().then(result => {
        console.log("Eliminado!")
        res.send("Eliminado!")
    }).catch(err => console.error(err))
});
```
- Obtener todos los registros de Mongo

Este endpoint es básico para poder graficar y llenar las tablas solicitadas con lo registrado en la base de datos.
```js
app.get('/obtenerUsuarios', async (req, res) => {
    coleccion.find().toArray()
    .then(results => {
        console.log("Obtener Usuarios!");
        res.json(results)
    })
    .catch(error => console.error(error))
});
```

### Dockerizando
Para colocar la Api y el servicio de base de datos MongoDB en un contenedor cada uno, se creó un archivo Dockerfile para la Api de NodeJs y se bajó una imagen pública de Mongo para la base de datos.

#### Dockerfile

Se especifica que el servicio es de Node y se obtiene la última versión. Se especifica el directorio de trabajo para NodeJs. Se copian los dos archivos JSON donde están especificadas todas las dependencias necesarias para la API a la raíz del contenedor. Se corre el comando ```npm install``` para que instale las dependencias. Se copian todos los archivos de raíz a raíz del contenedor. Se expone el puerto 8080, que es donde se va a exponer la API. Se corre el comando ```mkdir -p /elements/procs``` para crear la carpeta dentro del contenedor y así ahí poder montar el módulo de procesos y el módulo de RAM. Por último, se levanta la Api con ```node index.js```.
```dockerfile
FROM node:latest
WORKDIR /usr/src/nodejs
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 8080
RUN mkdir -p /elements/procs

CMD ["node", "index.js"]
```
#### Docker Compose
Se levantan los dos contenedores con Docker Compose, en el archivo se crea la red ```networkapi``` y se define cada contenedor con su nombre, puertos que cada contenedor va a exponer, su ruta donde puede encontrar el Dockerfile o ya sea el nombre de la imagen a bajar para el caso de Mongo.

```yml
version: "2.2"
services:  
  apinode:
    container_name: apinode
    restart: always
    build: ./nodejs
    ports:
      - "8080:8080"
    links:
      - db
    networks:
      - networkapi
    volumes:
      - /proc/:/elements/procs/

  db:
      image: 'mongo'
      container_name: db
      environment:
          - PUID=1000
          - PGID=1000
      volumes:
          - /home/barry/db/database:/data/db
      ports:
          - 27017:27017
      restart: unless-stopped
      networks:
        - networkapi

networks:
  networkapi:
    driver: "bridge"
```