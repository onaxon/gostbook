zeebo/gostbook clone to Learn Go

See http://web.archive.org/web/20130813221848/http://shadynasty.biz/blog/2012/07/30/quick-and-clean-in-go/

### Interesting bits
- Authenticate users using gorilla(sessions) and the coockie store
    + Register users
    + Login / logout
- Contex management via a custom handler
- Persist data in MongoDB
- Templates managed through a function

### ToDo
- Replace custom Context with gorilla/context
- Modularize with packages
- Reimplement with Martini? (search for martini tutorial)