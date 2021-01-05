Seneca
---

https://en.wikipedia.org/wiki/Seneca_the_Younger

This uses sci-hub, which is probably illegal.

Literature indexing, search, and note management. Special attention tended to: 

  * decentralized structure
  * low mental transaction cost on note read/write
  * speed
  * terminal based workflow with hackable internals

Assumes text editor of choice is `vim`. Good synergy with tiling window managers.


_Think of big soup of literature, don't impose structure_
1. Add paper with just URL
2. Search across text easily
3. Manage note buffers with each paper

`seneca <paper url>`
`seneca letters <grep search>`

-> Search just brings `head -l 5` - can easily identify title.
Then open paper with editor and open `seneca letter <paper>` for note.

## Implementation

Efficiently searchable directory structure
Directory for each paper:
  * pdf
  * pdf txt
  * note buffer

`seneca http` - pulls from scihub?

## Debt

  * Doesn't work on:
    * arxiv
    * ncbi

## Backlog

  * Better metadata parsing / terminal presentation
  * Grep over _only_ abstract or body 
  * Integration with Athena for regular (morning) reading across topic space
  * Automatic collection from arxiv based on things I should know about / read
