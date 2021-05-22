[TOC]

## 1  Two Sum

```java
class Solution {
    public int[] twoSum(int[] nums, int target) {
        HashMap<Integer,Integer> map = new HashMap<>();//(target - nums[i], i)
        int len = nums.length;
        for (int i = 0; i < len; i++) {
            if(map.containsKey(nums[i])) {
                return new int[] {map.get(nums[i]), i};
            } else {
                map.put(target - nums[i], i);
            }
        }
        return null;
    }
}
```

```python
class Solution(object):
    def twoSum(self, nums, target):
        """
        :type nums: List[int]
        :type target: int
        :rtype: List[int]
        """
        # map(nums[i], i) 若果target-nums[i]存在，则返回
        length = len(nums)
        numMap = dict()
        for i in range(0, length):
            targetNum = target - nums[i]
            if targetNum in numMap:
                return [numMap[targetNum], i]
            else:
                numMap[nums[i]] = i
        return None
```

```go
func twoSum(nums []int, target int) []int {
	// map(nums[i], i) return i, map.get(target - nums[i])
	freq := make(map[int]int)
	//freq := map[int]int{}
	for i, num := range nums {
		if j, ok := freq[target - num]; ok {
			return []int{j, i}
		} else {
			freq[num] = i
		}
	}
    return nil
}
```

## 2  Add Two Numbers

```java
class Solution {
    public ListNode addTwoNumbers(ListNode l1, ListNode l2) {
        ListNode head = null;
        ListNode tail = null;
        int carry = 0;
        while (l1 != null || l2 != null || carry != 0) {
            int sum = carry;
            if(l1 != null) {
                sum += l1.val;
                l1 = l1.next;
            }
            if(l2 != null) {
                sum += l2.val;
                l2 = l2.next;
            }
            carry = sum / 10;
            ListNode newNode = new ListNode(sum % 10);
            if(head == null) {
                head = newNode;
                tail = head;
            } else {
                tail.next = newNode;
                tail = tail.next;
            }
        }
        return head;
    }
}
```

```python
class Solution(object):
    def addTwoNumbers(self, l1, l2):
        """
        :type l1: ListNode
        :type l2: ListNode
        :rtype: ListNode
        """
        head = None
        tail = None
        carry = 0
        while l1 is not None or l2 is not None or carry != 0:
            total = carry
            if l1 is not None:
                total += l1.val
                l1 = l1.next
            if l2 is not None:
                total += l2.val
                l2 = l2.next
            carry = total / 10
            total = total % 10
            newNode = ListNode(total)
            if head is None:
                head = newNode
                tail = head
            else:
                tail.next = newNode
                tail = tail.next
        return head
```

```go
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
   // 思路：当l1或l2或有进位时，一直循环遍历
   var head *ListNode
   var tail *ListNode
   var carry int
   for l1 != nil || l2 != nil || carry != 0 {
      var sum = carry
      if l1 != nil {
         sum += l1.Val
         l1 = l1.Next
      }
      if l2 != nil {
         sum += l2.Val
         l2 = l2.Next
      }
      carry = sum / 10
      sum = sum % 10
      nextNode := &ListNode{
         sum,
         nil,
      }
      if head == nil {
         head = nextNode
         tail = head
      } else {
         tail.Next = nextNode
         tail = tail.Next
      }
   }
    return head
}
```

## 3  Longest Substring Without Repeating Characters

```java
class Solution {
    public int lengthOfLongestSubstring(String s) {
        int len = s.length();
        int maxLen = 0;
        int left = -1;//(left, right]区间内保存不重复字符串
        HashMap<Character, Integer> map = new HashMap<>();
        for (int right = 0; right < len; right++) {
            char c = s.charAt(right);
            if(map.containsKey(c) && map.get(c) > left) {
                left = map.get(c);
            }
            if(right - left > maxLen) {
                maxLen = right - left;
            }
            map.put(c, right);
        }
        return maxLen;
    }
}
```

## 46  Permutations

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
