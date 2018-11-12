# Showcase

On startup, the server web scrapes https://unity.com/madewith to grab all of the projects on the website to show to the client. The information is gathered into an array.
When a client requests the server, the server will return an HTML response of a random project within the array. This is implemented with a map from a user's remote address to a list of available projects to show.


Explanation:
When a new user pings the server, a new entry is made for them in the map with a key of their remote address to a new list of indices for the projects array. Then, we choose a random index in the array and return the project corresponding to that index. We then move that index to the end of the list, and "decrease" the size of our array by 1.
When the "size" of an array reaches 0, we simply change it back to its original size and continue choosing a random element as usual.

Because the chosen index is moved to the back of the array, it is guaranteed it will not be chosen until the array resets (when the size hits 0). The map also guarantees uniqueness for each user that pings the server. Because each integer array corresponds to exactly one remote address, one user will not see a repeated project until they have seen all of the other projects.

I chose this approach because all calls to the server would be O(1) runtime except for O(n) (where n is the number of projects) when a new user pings the server. This is therefore an amortized O(1) runtime.


In the future:
I would store the projects array and the user map within a database. Because all of this information is stored within the one go file, it is only accessible within that running instance of the program. Other servers of the same code base will not be able to share each other's maps. So, if we had multiple servers and the user pings different ones, it is certainly not guaranteed they wouldn't see repeated projects. Also, if this one server went down, all of map data would be gone and we wouldn't be able to retrieve it. Therefore, a database would help this server be persistent and elastic.

I would not use web scraping as my reliable source of gathering information. Because I have no access to any databases within Unity, the only way I could gather the information was through web scraping. If I were to build this as an actual product, I would request access to the database or an API that would give me the projects in a more reliable manner.
