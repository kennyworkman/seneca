Seneca 🏛️
---

https://en.wikipedia.org/wiki/Seneca_the_Younger

This uses sci-hub, which is probably illegal.

Literature indexing, search, and note management.

![](attention.gif)

### Quickstart

To install:

```
git clone https://github.com/kennyworkman/seneca
brew install zathura
make build
```

To add a paper:

```
seneca https://www.nature.com/articles/s41467-020-18008-4
```

To search across indexed papers and retrieve note buffer:

```
seneca l 
# or
seneca letters
```

I think you will find this is all that you need. There is no delete. If you
spent enough time indexing a reference, I'm of the opinion you should hang on to
it (just in case).

Notes persist across time. Opening a long forgotten paper weeks later will
recover time stamped notes from my last interaction with it.

### Design

 Special attention towards: 

  * decentralized structure
  * low mental transaction cost on note read/write
  * speed
  * terminal based workflow with hackable internals

Assumes text editor of choice is `vim`. Good synergy with tiling window managers.

Read more [here](https://kennethworkman.com/code/seneca/).

## Debt

  * Non-filesystem persistence - sqlite3
  * Support for non DOI-based references:
    * market research
    * books
    * technical preprints
  * Integration with other tools:
    * Automatic indexing of relevant papers at spaced interval
    * Priority list / backlog of desired reading
    * Automatic presentation of backlog depending on calendar events
