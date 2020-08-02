# 🌿gbr

Switch git branch interactively using arrow keys

## Demo

![](demo.gif)

## How to use

### Requirements

You need to have `node` installed.

### Install

```bash
yarn global add git-gbr
```

OR

```bash
npm install --global git-gbr
```

#### Maybe

If you're using zsh or any other kind of shell that already aliased "gbr" to
some other command, you may need to override this in your `~/.zshrc` /
`~/.bashrc` or similar file.

```bash
alias gbr='/usr/local/bin/gbr'
```

### Run

Now in any git repository, run

```bash
gbr
```
