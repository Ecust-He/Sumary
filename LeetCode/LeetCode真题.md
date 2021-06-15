[TOC]

## 1  Two Sum

### 思路

- 使用Map数据结构

### Java

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

### Python

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

### Go

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

### 思路

- 需要考虑进位

### Java

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

### Python

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

### Go

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

### 思路

- 首先需要定义[left, right]区间内保存不重复字符串
- 使用map，key为当前待考察的字符，value为索引位置

### Java

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

### Python

```python
class Solution(object):
    def lengthOfLongestSubstring(self, s):
        """
        :type s: str
        :rtype: int
        """
        maxLen = 0
        left = -1
        charMap = dict()
        for right in range(0, len(s)):
            cur = s[right]
            if cur in charMap and charMap[cur] > left:
                left = charMap[cur]
            if right - left > maxLen:
                maxLen = right - left
            charMap[cur] = right
        return maxLen
```

### Go

```go
func lengthOfLongestSubstring(s string) int {
   var left = -1//(left, right]区间内不重复
   var charToIndexMap = make(map[rune]int)
   var maxLength int
   for right, c := range []rune(s) {
      if preIndex, ok := charToIndexMap[c]; ok && preIndex > left {
         left = preIndex
      }
      if right - left > maxLength {
         maxLength = right - left
      }
      charToIndexMap[c] = right
   }
    return maxLength
}
```

## 4  Median of Two Sorted Arrays

### 思路

- 需要开辟新的空间保存合并之后的数组

### Java

```java
public static double findMedianSortedArrays(int[] nums1, int[] nums2) {
    int m = nums1.length;
    int n = nums2.length;
    int i = 0;
    int j = 0;
    int k = 0;
    int[] num = new int[m + n];// 合并之后的有序数组
    while (i < m || j < n) {
        if(i < m && j >= n) {
            num[k++] = nums1[i++];
        } else if(j < n && i >= m) {
            num[k++] = nums2[j++];
        } else {
            if (nums1[i] <= nums2[j]) {
                num[k++] = nums1[i++];
            } else {
                num[k++] = nums2[j++];
            }
        }
    }
    double res = 0;
    if ((m + n) % 2 == 1) {
        res = num[(m + n) / 2];
    } else {
        res = ((double) num[(m + n) / 2] + (double) num[(m + n) / 2 -1]) / 2;
    }
    return res;
}
```

### Python

```python
def findMedianSortedArrays(self, nums1, nums2):
    """
    :type nums1: List[int]
    :type nums2: List[int]
    :rtype: float
    """
    # 使用切片
    num = list()
    while len(nums1) > 0 or len(nums2) > 0:
        if len(nums1) > 0 and len(nums2) > 0:
            if nums1[0] < nums2[0]:
                num.append(nums1[0])
                nums1 = nums1[1:]
            else:
                num.append(nums2[0])
                nums2 = nums2[1:]
            continue
        if len(nums1) > 0:
            num.append(nums1[0])
            nums1 = nums1[1:]
        if len(nums2) > 0:
            num.append(nums2[0])
            nums2 = nums2[1:]
    length = len(num)
    if length % 2 == 0:
        return float(num[length/2] + num[length/2 - 1]) / 2
    else:
        return float(num[length/2])
```

### Go

```go
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
   var nums []int
   for len(nums1) != 0 || len(nums2) != 0 {
      if len(nums1) != 0 && len(nums2) != 0 {
         if nums1[0] < nums2[0] {
            appendNum(&nums1, &nums)
         } else {
            appendNum(&nums2, &nums)
         }
         continue
      }
      if len(nums1) != 0 {
         appendNum(&nums1, &nums)
      }
      if len(nums2) != 0 {
         appendNum(&nums2, &nums)
      }
   }
   var length = len(nums)
   if length%2 == 0 {
      return float64(nums[length/2]+nums[length/2-1]) / 2
   } else {
      return float64(nums[length/2])
   }
}
```

## 7  Reverse Integer

### 思路

- 拨数

### Java

```java
public int reverse(int x) {
    int ret = 0;
    while (x != 0) {
        ret = ret * 10 + x % 10;
        if (ret < Integer.MIN_VALUE || ret > Integer.MAX_VALUE) return 0;
        x /= 10;
    }
    return ret;
}
```

### Python

```python
class Solution(object):
    def reverse(self, x):
        """
        :type x: int
        :rtype: int
        """
        y = 0
        while x != 0:
            y = y * 10 + x % 10
            if y > int(pow(2, 31) - 1) or y < int(pow(-2, 31)):
                return 0
            x /= 10
        return y
```

### Go

```java
func reverse(x int) int {
   var ret int
   max_int := int(math.Pow(2, 31) - 1)
   min_int := int(math.Pow(-2, 31))
   for x != 0 {
      ret = ret * 10 + x % 10
      if ret > max_int || ret < min_int {
         return 0
      }
      x /= 10
   }
    return ret
}
```

## 9  Palindrome Number

### 思路

- 双指针对撞

### Java

```java
class Solution {
    public boolean isPalindrome(int x) {
        if (x < 0) {
            return false;
        }
        String str = String.valueOf(x);
        int i = 0;
        int j = str.length() - 1;
        while (i < j) {
            if(str.charAt(i++) != str.charAt(j--)) {
                return false;
            }
        }
        return true;
    }
}
```

### Python

```python
class Solution(object):
    def isPalindrome(self, x):
        """
        :type x: int
        :rtype: bool
        """
        if x < 0:
            return False
        y = str(x)
        # i = 0
        # j = len(y) - 1
        # while i < j:
        #     if y[i] != y[j]:
        #         return False
        #     i += 1
        #     j -= 1
        while y:
            if y[0] != y[-1]:
                return False
            y = y[1:-1]
        return True
```

### Go

```go
func isPalindrome(x int) bool {
   if x < 0 {
      return false
   }
   y := strconv.Itoa(x)
   i := 0
   j := len(y) - 1
   for i < j {
      if y[i] != y[j] {
         return false
      }
      i += 1
      j -= 1
   }

   //for len(y) > 1 {
   // if y[0] != y[len(y)-1] {
   //    return false
   // }
   // y = y[1:len(y)-1]
   //}
   return true
}
```

## 11  container-with-most-water

### 思路

已短板一边为准

### Java

```java
class Solution {
    public int maxArea(int[] height) {
        int i = 0;
        int j = height.length - 1;
        int maxArea = 0;
        while (i < j) {
            if(height[i] < height[j]) {
                maxArea = Math.max(maxArea, (j - i) * height[i]);
                i++;
            } else {
                maxArea = Math.max(maxArea, (j - i) * height[j]);
                j--;
            }
        }
        return maxArea;
    }
}
```

### Python

```python
class Solution(object):
    def maxArea(self, height):
        """
        :type height: List[int]
        :rtype: int
        """
        left = 0
        right = len(height) - 1
        maxArea = 0
        while left < right:
            if height[left] < height[right]:
                minHeight = height[left]
                left += 1
            else:
                minHeight = height[right]
                right -= 1
            area = minHeight * (right - left + 1)
            if area > maxArea:
                maxArea = area
        return maxArea
```

### Go

```go
func maxArea(height []int) int {
   var left = 0
   var right = len(height) - 1
   var maxArea = 0
   for left < right {
      var area int
      if height[left] < height[right] {
         area = height[left] * (right - left)
         left++
      } else {
         area = height[right] * (right - left)
         right--
      }
      if area > maxArea {
         maxArea = area
      }
   }
   return maxArea
}
```

## 13  Roman to Integer

### 思路

- 先穷举出所有情况
- 先看下能否两位取出再考虑一位取出解析

### Java

```java
class Solution {
    public int romanToInt(String s) {
        HashMap<String, Integer> symbolMap = new HashMap();
        symbolMap.put("I", 1);
        symbolMap.put("IV", 4);
        symbolMap.put("V", 5);
        symbolMap.put("IX", 9);
        symbolMap.put("X", 10);
        symbolMap.put("XL", 40);
        symbolMap.put("L", 50);
        symbolMap.put("XC", 90);
        symbolMap.put("C", 100);
        symbolMap.put("CD", 400);
        symbolMap.put("D", 500);
        symbolMap.put("CM", 900);
        symbolMap.put("M", 1000);

        int sum = 0;
        int index = 0;
        int len = s.length();
        while (index < len) {
            if (len - index > 1 && symbolMap.containsKey(s.substring(index, index + 2))) {
                sum += symbolMap.get(s.substring(index, index + 2));
                index += 2;
            } else {
                sum += symbolMap.get(s.substring(index, index + 1));
                index += 1;
            }
        }
        return sum;
    }
}
```

### Python

```python
class Solution(object):
    def romanToInt(self, s):
        """
        :type s: str
        :rtype: int
        """
        romanToIntDict = {
            'I': 1,
            'IV': 4,
            'V': 5,
            'IX': 9,
            'X': 10,
            'XL': 40,
            'L': 50,
            'XC': 90,
            'C': 100,
            'CD': 400,
            'D': 500,
            'CM': 900,
            'M': 1000
        }

        total = 0
        while len(s):
            if len(s) > 1 and s[0: 2] in romanToIntDict:
                total += romanToIntDict[s[0: 2]]
                s = s[2:]
            else:
                total += romanToIntDict[s[0: 1]]
                s = s[1:]
        return total
```

### Go

```go
func romanToInt(s string) int {
   // 思路：如果字符串长度大于2，优先判断map中的key是否存在，若存在则取出来
   // 若不存在，先解析1位,如此反复
   romanToIntMap := map[string]int{
      "I": 1,
      "IV": 4,
      "V": 5,
      "IX": 9,
      "X": 10,
      "XL": 40,
      "L": 50,
      "XC": 90,
      "C": 100,
      "CD": 400,
      "D": 500,
      "CM": 900,
      "M": 1000,
   }
   var sum int
   for len(s) > 0 {
        if len(s) > 1 {
            if num, ok := romanToIntMap[s[0: 2]];ok {
                sum += num
                s = s[2:]
                continue
            }
        }
      if num, ok := romanToIntMap[string(s[0])];ok {
         sum += num
         s = s[1:]
      }

    }
    return sum
}
```

## 14  Longest Common Prefix

### 思路

- 选取数组中第一个元素作为LCP
- 依次与数组中其他元素做交集处理

### Java

```java
class Solution {
    public String longestCommonPrefix(String[] strs) {
        int len = strs.length;
        if(len == 0) {
            return "";
        }
        if(len == 1) {
            return strs[0];
        }
        String commonPrefix = strs[0];
        for (int i = 1; i < len; i++) {
            while (!strs[i].startsWith(commonPrefix)) {
                if (commonPrefix.length() > 1) {
                    commonPrefix = commonPrefix.substring(0, commonPrefix.length() - 1);
                } else {
                    return "";
                }
            }
        }
        return commonPrefix;
    }
}
```

### Python

```python
def longestCommonPrefix(self, strs):
    """
    :type strs: List[str]
    :rtype: str
    """
    if len(strs) == 0:
        return ''
    if len(strs) == 1:
        return strs[0]
    commonPrefix = strs[0]
    for i in range(1, len(strs)):
        while not strs[i].startswith(commonPrefix):
            if len(commonPrefix) > 0:
                commonPrefix = commonPrefix[:len(commonPrefix) - 1]
            else:
                return ''
    return commonPrefix
```

### Go

```go
func longestCommonPrefix(strs []string) string {
   if len(strs) == 0 {
      return ""
   }
   if len(strs) == 1 {
      return strs[0]
   }
   commonPrefix := strs[0]
   for _, str := range strs {
      for !strings.HasPrefix(str, commonPrefix) {
         if len(commonPrefix) >= 1 {
            commonPrefix = commonPrefix[:len(commonPrefix)-1]
         } else {
            return ""
         }
      }
   }
   return commonPrefix
}
```

## 15  3Sum

### 思路

### Java

### Python

### Go

## 17  Letter Combinations of a Phone Number

### 思路

### Java

```java
class Solution {
    private String[] digitMap = {
            "",
            "",
            "abc",
            "def",
            "ghi",
            "jkl",
            "mno",
            "pqrs",
            "tuv",
            "wxyz"
    };
    public List<String> letterCombinations(String digits) {
        List<String> res = new ArrayList<>();
        if(digits.isEmpty()) {
            return res;
        }
        // [0, index)解析的字符保存在str中
        letterCombinations(0, "", digits, res);
        return res;
    }

    public void letterCombinations(int index, String str, String digits, List<String> res) {
        if (str.length() == digits.length()) {
            res.add(str);
            return;
        }
        String target = digitMap[digits.charAt(index) - '0'];
        for (int i = 0; i < target.length(); i++) {
            letterCombinations(index + 1, str + target.charAt(i), digits, res);
        }
    }
}
```

### Python

```python
# leetcode submit region begin(Prohibit modification and deletion)
class Solution(object):
    digitMap = [
            "",
            "",
            "abc",
            "def",
            "ghi",
            "jkl",
            "mno",
            "pqrs",
            "tuv",
            "wxyz"
    ]

    def letterCombinations(self, digits):
        """
        :type digits: str
        :rtype: List[str]
        """
        if not digits:
            return []
        res = []
        self.combinations(0, '', digits, res)
        return res

    def combinations(self, index, target_str, digits, res):
        """
        :param index: int
        :param target_str: str
        :param digits: str
        :param res: List[str]
        :return:
        """
        if len(target_str) == len(digits):
            res.append(target_str)
            return
        for ch in self.digitMap[int(digits[index])]:
            self.combinations(index + 1, target_str + ch, digits, res)
```

### Go

```go
var digitMap = []string{
    "",
    "",
    "abc",
    "def",
    "ghi",
    "jkl",
    "mno",
    "pqrs",
    "tuv",
    "wxyz",
}

func letterCombinations(digits string) []string {
    var res []string
    if len(digits) == 0 {
        return res
    }
    letterCombination(0, "", digits, &res)
    return res
}

/**
  digits[0, index)内解析后的字符串保存在str
*/
func letterCombination(index int, str string, digits string, res *[]string) {
    if len(str) == len(digits) {
        *res = append(*res, str)
        return
    }
    for _, ch := range digitMap[digits[index] -48] {
        letterCombination(index + 1, str + string(ch), digits, res)
    }
}
```

## 19  Remove Nth Node From End if List

### 思路

### Java

```java
class Solution {
    public ListNode removeNthFromEnd(ListNode head, int n) {
        int sz = 0;
        ListNode cur = head;
        while (cur != null) {
            sz++;
            cur = cur.next;
        }
        ListNode dumyHead = new ListNode(0);
        dumyHead.next = head;
        ListNode pre = dumyHead;
        for (int i = 0; i < sz - n; i++) {
            pre = pre.next;
        }
        pre.next = pre.next.next;
        return dumyHead.next;
    }
}
```

### Python

### Go

## 46  Permutations

### Java

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
