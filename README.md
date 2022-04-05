# GO-MOVIES

Go-Movies is a movie-catalogue management application where there are two 
broad categories of users. Authenticated Users and unathenticated users.

Go-movies was built using the GO programming language as well as other 
technologies such as: 
- **GRAPHQL API** for the search feature
- **REST API** for CRUD operations
- **JWT** for authentication
- **BCRYPT** for encryption
- **POSTGRES** to store data


### PRIVILEGES FOR UNAUTHENTICATED USERS

- Users can view all the movies in the catalogue
- Users can view details about a particular movie in the catalogue
- Users can search by title for a movie they are curious about
- Users can view movies by their genres

### PRIVILEGES FOR AUTHENTICATED USERS

Users have all the priviliges that unauthenticated users have including the following:

- Users can add a new movie to the catalogue
- Users can edit a previously entered movie in the catalogue
- Users can delete a movie from the catalogue