[TOC]

## 46	Permutations

```java
class Solution {
    private ArrayList<List<Integer>> res;
    private boolean[] hasSelected;

    public List<List<Integer>> permute(int[] nums) {
        if(nums == null || nums.length == 0) {
            return res;
        }
        res = new ArrayList<List<Integer>>();
        hasSelected = new boolean[nums.length];
        permute(nums, new LinkedList<Integer>());
        return res;
    }

    //	每次从数组nums中选取一个不重复的数，放入集合selectList
    public void permute(int[] nums, LinkedList<Integer> selectedList) {
        if(selectedList.size() == nums.length) {
            res.add(new ArrayList<Integer>(selectedList));
            return;
        }
        for (int i = 0; i < nums.length; i++) {
            if(!hasSelected[i]) {
                hasSelected[i] = true;
                selectedList.add(nums[i]);
                permute(nums, selectedList);
                selectedList.removeLast();
                hasSelected[i] = false;
            }
        }
        return;
    }
}
```