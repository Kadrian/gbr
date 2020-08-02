import { execSync } from "child_process";
import { prompt } from "inquirer";

function getRecentBranchesWithDates() {
  const branchesWithDateCommand =
    "git for-each-ref --sort=committerdate refs/heads/ --format='%(committerdate:short) %(refname:short)'";

  return execSync(branchesWithDateCommand)
    .toString()
    .trim()
    .split("\n")
    .map((br) => br.trim());
}

function getCurrentBranchRef() {
  return execSync("git branch --show-current").toString().trim();
}

function getCurrentBranchIdx(branches: string[]) {
  const branchRefs = branches.map((br) => br.split(" ")[1]);
  const currentBranchRef = getCurrentBranchRef();

  return branchRefs.indexOf(branchRefs.find((br) => br === currentBranchRef));
}

async function switchBranch() {
  const branches = getRecentBranchesWithDates();
  const currentBranchIdx = getCurrentBranchIdx(branches);

  const choice = await prompt({
    name: "Branch",
    type: "list",
    default: currentBranchIdx,
    choices: branches,
    loop: false,
  });

  const nextBranch = choice.Branch.split(" ")[1];

  execSync(`git checkout ${nextBranch}`);
}

(async function () {
  try {
    await switchBranch();
  } catch (e) {
    // All errors in os commands are printed already, so we can simply catch and exit(1)
    process.exit(1);
  }
})();
