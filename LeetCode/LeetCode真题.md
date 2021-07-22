[TOC]

## 1  ~~Two Sum~~

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

## 2  ~~Add Two Numbers~~

### 思路

- 加法进位

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
- 使用map数据结构，key为当前待考察的字符，value为索引位置

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

## 4  ~~Median of Two Sorted Arrays~~

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

## 11  ~~container-with-most-water~~

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

## 14  ~~Longest Common Prefix~~

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

## 19  ~~Remove Nth Node From End if List~~

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

```python
class Solution(object):
    def removeNthFromEnd(self, head, n):
        """
        :type head: ListNode
        :type n: int
        :rtype: ListNode
        """
        sz = 0
        curNode = head
        while curNode:
            sz += 1
            curNode = curNode.next
        dumyHead = ListNode(0)
        dumyHead.next = head
        preNode = dumyHead
        for i in range(0, sz - n):
            preNode = preNode.next
        preNode.next = preNode.next.next
        return dumyHead.next
```

### Go

```go
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    n := 0
    curNode := head
	for curNode != nil {
		n += 1
		curNode = curNode.Next
	}
	dumyHead := new(ListNode)
	dumyHead.Next = head
    preNode := dumyHead
	for i := 0; i < sz - n; i++ {
		preNode = preNode.Next
	}
	preNode.Next = preNode.Next.Next
    return dumyHead.Next
}
```

## 20  ~~Valid Parentheses~~

### 思路

构造栈数据结构

### Java

```java
public boolean isValid(String s) {
    Stack<Character> stack = new Stack<Character>();
    int len = s.length();
    for (int i = 0; i < len; i++) {
        Character c = s.charAt(i);
        if(c == '{' || c == '[' || c == '(') {
            stack.push(c);
        } else {
            if (stack.isEmpty()) {
                return false;
            }
            Character top = stack.pop();
            if((top == '{' && c != '}') || (top == '[' && c != ']') || (top == '(' && c != ')')) {
                return false;
            }
        }
    }
    if(!stack.isEmpty()) {
        return false;
    }
    return true;
}
```

### Python

```python
def isValid(self, s):
    """
    :type s: str
    :rtype: bool
    """
    stack = list()
    for c in s:
        if c in ('{', '[', '('):
            stack.append(c)
        else:
            if len(stack) < 1:
                return False
            top = stack.pop()
            if (c == '}' and top != '{') or (c == ']' and top != '[') or (c == ')' and top != '('):
                return False
    if len(stack) > 0:
        return False
    return True
```

### Go

```go
func isValid(s string) bool {
    var stack Stack
    for _, byte := range s {
        c := string(byte)
        if c == "(" || c == "{" || c == "[" {
            stack.push(c)
        } else {
            if stack.isEmpty() {
                return false
            }
           head := stack.pop()
            if (c == ")" && head != "(") || (c == "}" && head != "{") || (c == "]" && head != "[") {
                return false
            }
        }
    }
    if stack.isNotEmpty() {
        return false
    }
    return true
}

type Stack []string

func (stack *Stack) push(s string) {
    *stack = append(*stack, s)
}

func (stack *Stack) pop() string {
    top := (*stack)[len(*stack)-1]
    *stack = (*stack)[:len(*stack)-1]
    return top
}

func (stack Stack) isEmpty() bool {
    return len(stack) <= 0
}

func (stack Stack) isNotEmpty() bool {
    return len(stack) > 0
}
```

## 21  ~~Merge Two Sorted Lists~~

### Java

```java
class Solution {
    public ListNode mergeTwoLists(ListNode l1, ListNode l2) {
        List<Integer> list = new ArrayList<>();
        while (l1 != null || l2 != null) {
            if(l1 != null && l2 != null) {
                if (l1.val <= l2.val) {
                    list.add(l1.val);
                    l1 = l1.next;
                } else {
                    list.add(l2.val);
                    l2 = l2.next;
                }
            } else if (l1 != null) {
                list.add(l1.val);
                l1 = l1.next;
            } else {
                list.add(l2.val);
                l2 = l2.next;
            }
        }
        ListNode dumyHead = new ListNode(0);
        ListNode cur = dumyHead;
        for (Integer num : list) {
            cur.next = new ListNode(num);
            cur = cur.next;
        }
        return dumyHead.next;
    }
}
```

### Python

```python
def mergeTwoLists(self, l1, l2):
    """
    :type l1: ListNode
    :type l2: ListNode
    :rtype: ListNode
    """
    res = []
    while l1 or l2:
        if l1 and l2:
            if l1.val <= l2.val:
                res.append(l1.val)
                l1 = l1.next
            else:
                res.append(l2.val)
                l2 = l2.next
        elif l1:
            res.append(l1.val)
            l1 = l1.next
        else:
            res.append(l2.val)
            l2 = l2.next
    dumyhead = ListNode(0)
    curNode = dumyhead
    for num in res:
        curNode.next = ListNode(num)
        curNode = curNode.next
    return dumyhead.next
```

### Go

```go
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
   var valSlice []int
   for l1 != nil || l2 != nil {
      if l1 != nil && l2 != nil {
         if l1.Val < l2.Val {
            valSlice = append(valSlice, l1.Val)
            l1 = l1.Next
         } else {
            valSlice = append(valSlice, l2.Val)
            l2 = l2.Next
         }
         continue
      }
      if l1 != nil {
         valSlice = append(valSlice, l1.Val)
         l1 = l1.Next
         continue
      }
      if l2 != nil {
         valSlice = append(valSlice, l2.Val)
         l2 = l2.Next
         continue
      }
   }
   dumyHead := new(ListNode)
   p := dumyHead
   for _, val := range valSlice {
      p.Next = &ListNode{val, nil}
      p = p.Next
   }
    return dumyHead.Next
}
```

## 24  Swap Nodes in Pairs

### 思路

需要四个变量preNode、node1、node2、nextNode

### Java

```java
    public ListNode swapPairs(ListNode head) {
        ListNode dumyHead = new ListNode(0, head);
        ListNode preNode = dumyHead;
        while(preNode.next != null && preNode.next.next != null){
            ListNode node1 = preNode.next;
            ListNode node2 = node1.next;
            ListNode next = node2.next;
            
            
            preNode.next = node2;
            node2.next = node1;
            node1.next = next;
            
            preNode = node1;
        }
        return dumyHead.next;
    }
```

### Python

```python
def swapPairs(self, head):
    """
    :type head: ListNode
    :rtype: ListNode
    """
    dumyHead = ListNode(0)
    dumyHead.next = head
    predode = dumyHead
    while predode.next and predode.next.next:
        node1 = predode.next
        node2 = node1.next
        nextNode = node2.next

        predode.next = node2
        node2.next = node1
        node1.next = nextNode

        predode = node1
    return dumyHead.next
```

### Go

```go
func swapPairs(head *ListNode) *ListNode {
    dumyHead := ListNode{Val:0, Next:head}
    preNode := &dumyHead
    for preNode.Next != nil && preNode.Next.Next != nil {
        node1 := preNode.Next
        node2 := node1.Next
        nextNode := node2.Next
        
        preNode.Next = node2
        node2.Next = node1
        node1.Next = nextNode
        
        preNode = node1
    }
    return dumyHead.Next;
}
```

## 26  ~~Remove Duplicates from Sorted Array~~

### Java

```java
public int removeDuplicates(int[] nums) {
    if(nums.length == 0){
        return 0;
    }
    int len = 1;//[0,len)保存不重复num
    for (int i = 1; i < nums.length; i++) {
        if(nums[i] != nums[len - 1]){
            nums[len] = nums[i];
            len++;
        }
    }
    return len;
}
```

### Python

```python
def removeDuplicates(self, nums):
    """
    :type nums: List[int]
    :rtype: int
    """
    if len(nums) == 0:
        return 0
    # [0, n)区间内保存不重复num
    n = 1
    for i in range(1, len(nums)):
        if nums[i] != nums[n-1]:
            nums[n] = nums[i]
            n += 1
    return n
```

### Go

```go
func removeDuplicates(nums []int) int {
   if len(nums) == 0{
      return 0
   }
   n := 1// [0,n)区间内保存不重复num
   for i := 1;i < len(nums);i++{
      if nums[i] != nums[n-1]{
         nums[n] = nums[i]
         n += 1    
      }
   }
   return n
}
```

## 27  ~~Remove Element~~

### Java

```java
public int removeElement(int[] nums, int val) {
    int len = 0;//[0, len)区间内不包含val
    for(int i = 0; i < nums.length; i ++){
        if(nums[i] != val) {
            nums[len++] = nums[i];
        }            
    }    
    return len;
}
```

### Python

```python
def removeElement(self, nums, val):
    """
    :type nums: List[int]
    :type val: int
    :rtype: int
    """
    n = 0
    for num in nums:
        if num != val:
            nums[n] = num
            n += 1
    return s
```

### Go

```go
func removeElement(nums []int, val int) int {
    n := 0
    for i := 0; i < len(nums); i++ {
        if(nums[i] != val) {
            nums[n] = nums[i]
            n += 1
        }
    }
    return n
}
```

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

### Python

### Go

## 53   ~~Maximum Subarray~~

### 思路

动态规划思想

### Java

```java
    public int maxSubArray(int[] nums) {
        int len = nums.length;
        if(len ==0) {
            return 0;
        }
        if(len == 1){
            return nums[0];
        }
        int[] memo = new int[len];// memo[i]记录以nums[i]结尾的子数组和的最大值
        memo[0] = nums[0];
        int max = memo[0];
        for(int i = 1; i < len; i++) {
            memo[i] = nums[i];
            if(memo[i-1] > 0) {
                memo[i] += memo[i-1];
            }
            if(memo[i] > max){
                max = memo[i];
            }
        }
        return max;
    }
```

### Python

```python
def maxSubArray(nums):
    """
    :type nums: List[int]
    :rtype: int
    """
    n = len(nums)
    if n == 0:
        return 0
    if n == 1:
        return nums[0]
    memo = [num for num in nums]  # 列表推导式
    for i in range(1, n):
        if memo[i - 1] > 0:
            memo[i] += memo[i - 1]
    return max(memo)
```

### Go

```go
func maxSubArray(nums []int) int {
    n := len(nums)
    if n == 0 {
        return 0
    }
    if n == 1 {
        return nums[0]
    }
    
    var memo = make([]int, n)
    memo[0] = nums[0]
    res := memo[0]
    for i := 1; i < n; i++ {
        memo[i] = nums[i]
        if memo[i -1] > 0 {
            memo[i] += memo[i-1]
        }
        if memo[i] > res {
            res = memo[i]
        }
    }
    return res
    
}
```

## 61  ~~Rotate List~~

### 思路

构造虚拟头结点，找到head、tail、newHead、newTail

### Java

```java
public ListNode rotateRight(ListNode head, int k) {
    if(head == null) {
        return null;
    }
    int n = 0;
    ListNode dumyHead = new ListNode(0, head);
    ListNode tail = dumyHead;
    while (tail.next != null) {
        n++;
        tail = tail.next;
    }

    int m = k % n;
    if(m == 0) {
        return head;
    }

    ListNode newTail = dumyHead;
    for (int i = 0; i < n - m; i++) {
        newTail = newTail.next;
    }
    ListNode newHead = newTail.next;
    tail.next = dumyHead.next;
    newTail.next = null;
    return newHead;
}
```

### Python

```python
def rotateRight(self, head, k):
    """
    :type head: ListNode
    :type k: int
    :rtype: ListNode
    """
    if head is None:
        return None
    dumyHead = ListNode(0, head)
    tail = dumyHead
    n = 0
    while tail.next:
        n += 1
        tail = tail.next
    if n == 1:
        return head
    m = k % n
    if m == 0:
        return dumyHead.next
    newTail = dumyHead
    for i in range(0, n - m):
        newTail = newTail.next
    newHead = newTail.next
    tail.next = dumyHead.next
    newTail.next = None
    return newHead
```

### Go

```go
func rotateRight(head *ListNode, k int) *ListNode {
    if head == nil {
        return nil
    }
    dumyHead := ListNode{Val:0, Next:head}
    tail := &dumyHead
    n := 0
    for ;tail.Next != nil;{
        n++
        tail = tail.Next
    }
    m := k % n
    if m == 0 {
        return head
    }
    newTail := &dumyHead
    for i := 0; i < n -m; i++ {
        newTail = newTail.Next
    }
    newHead := newTail.Next
    tail.Next = dumyHead.Next
    newTail.Next = nil
    return newHead
}
```

## 62  ~~Unique Paths~~

### Java

```java
public int uniquePaths(int m, int n) {
    int[][] memo = new int[m][n];
    if(m == 1) {
        return 1;
    }
    for (int i = 0; i < m; i++) {
        memo[i][0] = 1;
    }
    for (int j = 0; j < n; j++) {
        memo[0][j] = 1;
    }

    for (int i = 1; i < m; i++) {
        for (int j = 1; j < n; j++) {
            memo[i][j] = memo[i-1][j] + memo[i][j-1];
        }
    }
    return memo[m-1][n-1];
}
```

### Python

#### 方法一

```python
def uniquePaths(self, m, n):
    """
    :type m: int
    :type n: int
    :rtype: int
    """
    if m == 1:
        return 1
    # memo = [[1 for i in range(0, n)] for j in range(0, m)]
    memo = [[1 for i in range(0, n)]]*m
    for i in range(1, m):
        for j in range(1, n):
            memo[i][j] = memo[i-1][j] + memo[i][j-1]
    return memo[m-1][n-1]
```

#### 方法二

```python
def uniquePaths(self, m, n):
    """
    :type m: int
    :type n: int
    :rtype: int
    """
    if m == 1:
        return 1
    memo = [1 for i in range(0, n)]
    for i in range(1, m):
        for j in range(1, n):
            memo[j] += memo[j-1]
    return memo[n-1]
```

### Go

#### 方法一

```go
// //空间复杂度O(m * n)
func uniquePaths(m int, n int) int {
    if m == 1 {
        return 1
    }
    memo := make([][]int, m)
    for i := 0; i < m; i++ {
        memo[i] = make([]int, n)
        memo[i][0] = 1
    }

    for j := 0; j < n; j++ {
        memo[0][j] = 1
    }

    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            memo[i][j] = memo[i][j-1] + memo[i-1][j]
        }
    }
    return memo[m-1][n-1]
}
```

#### 方法二

```go
//空间复杂度O(n)
func uniquePaths(m int, n int) int {
    if m == 1 {
        return 1
    }
    memo := make([]int, n)
    for j := 0; j < n; j++ {
        memo[j] = 1
    }
    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            memo[j] += memo[j-1]
        }
    }
    return memo[n-1]
}
```

## 63  ~~Unique Paths II~~

### Java

```java
public int uniquePathsWithObstacles(int[][] obstacleGrid) {
    int m = obstacleGrid.length;
    int n = obstacleGrid[0].length;

    int[][] memo = new int[m][n];

    if(obstacleGrid[0][0] == 1) {
        return 0;
    } else {
      memo[0][0] = 1;
    }

    for (int i = 1; i < m; i++) {
        memo[i][0] = memo[i-1][0];
        if(obstacleGrid[i][0] == 1) {
            memo[i][0] = 0;
        }
    }
    for (int j = 1; j < n; j++) {
        memo[0][j] = memo[0][j-1];
        if(obstacleGrid[0][j] == 1) {
            memo[0][j] = 0;
        }
    }
    for (int i = 1; i < m; i++) {
        for (int j = 1; j < n; j++) {
            if(obstacleGrid[i][j] == 1) {
                memo[i][j] = 0;
            } else {
                memo[i][j] = memo[i][j-1] + memo[i-1][j];
            }
        }
    }
    return memo[m-1][n-1];
}
```

### Python

```python
def uniquePathsWithObstacles(self, obstacleGrid):
    """
    :type obstacleGrid: List[List[int]]
    :rtype: int
    """
    m = len(obstacleGrid)
    n = len(obstacleGrid[0])
    memo = [[0 for i in range(0, n)] for j in range(0, m)]
    if obstacleGrid[0][0] == 1:
        return 0
    else:
        memo[0][0] = 1
    for i in range(1, m):
        memo[i][0] = memo[i-1][0]
        if obstacleGrid[i][0] == 1:
            memo[i][0] = 0
    for j in range(1, n):
        memo[0][j] = memo[0][j-1]
        if obstacleGrid[0][j] == 1:
            memo[0][j] = 0
    for i in range(1, m):
        for j in range(1, n):
            if obstacleGrid[i][j] == 1:
                memo[i][j] = 0
            else:
                memo[i][j] = memo[i][j-1] + memo[i-1][j]
    return memo[m-1][n-1]
```

### Go

```go
func uniquePathsWithObstacles(obstacleGrid [][]int) int {
    m := len(obstacleGrid)
    n := len(obstacleGrid[0])
    memo := make([][]int, m)
    for i := 0; i < m; i++ {
        memo[i] = make([]int, n)
    }
    if obstacleGrid[0][0] == 1 {
        return 0
    } else {
        memo[0][0] = 1
    }
    for i := 1; i < m; i++ {
        memo[i][0] = memo[i-1][0]
        if obstacleGrid[i][0] == 1 {
            memo[i][0] = 0
        }
    }
    for j := 1; j < n; j++ {
        memo[0][j] = memo[0][j-1]
        if obstacleGrid[0][j] == 1 {
            memo[0][j] = 0
         }
    }
    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            if obstacleGrid[i][j] == 1 {
                memo[i][j] = 0
            } else {
                memo[i][j] = memo[i-1][j] + memo[i][j-1]
            }
        }
    }
    return memo[m-1][n-1]
}
```

## 64  ~~Minimum Path Sum~~

### Java

```java
public int minPathSum(int[][] grid) {
    int m = grid.length;
    int n = grid[0].length;
    int[][] memo = new int[m][n];
    memo[0][0] = grid[0][0];
    for (int i = 1; i < m; i++) {
        memo[i][0] = memo[i-1][0] + grid[i][0];
    }
    for (int j = 1; j < n; j++) {
        memo[0][j] = memo[0][j-1] + grid[0][j];
    }
    for (int i = 1; i < m; i++) {
        for (int j = 1; j < n; j++) {
            memo[i][j] = grid[i][j] + Math.min(memo[i-1][j], memo[i][j-1]);
        }
    }
    return memo[m-1][n-1];
}
```

### Python

### Go

## 67  Add Binary

### Java

```java
public String addBinary(String a, String b) {
    Stack<Integer> stack1 = putStrIntoStack(a);
    Stack<Integer> stack2 = putStrIntoStack(b);
    Stack<Integer> stack = new Stack<>();
    int carry = 0;
    while (!stack1.isEmpty() || !stack2.isEmpty() || carry != 0) {
        int sum = carry;
        if (!stack1.isEmpty()) {
            sum += stack1.pop();
        }
        if (!stack2.isEmpty()) {
            sum += stack2.pop();
        }
        carry = sum / 2;
        stack.add(sum % 2);
    }
    StringBuilder res = new StringBuilder();
    while (!stack.isEmpty()) {
        res.append(stack.pop());
    }
    return res.toString();
}

public Stack<Integer> putStrIntoStack(String s){
    Stack<Integer> stack = new Stack<>();
    for (int i = 0; i < s.length(); i++) {
        stack.add(Integer.valueOf(String.valueOf(s.charAt(i))));
    }
    return stack;
}
```

### Python

```python
def addBinary(self, a, b):
    """
    :type a: str
    :type b: str
    :rtype: str
    """
    stack1 = [s for s in a]
    stack2 = [s for s in b]
    stack = []
    carry = 0
    while stack1 or stack2 or carry:
        total = carry
        if stack1:
            total += int(stack1.pop())
        if stack2:
            total += int(stack2.pop())
        carry = total / 2
        total %= 2
        stack.append(total)
    res = ''    
    while stack:
        res += str(stack.pop())
    return res  
```

### Go

## 70  Climbing Stairs

### Java

```java
public int climbStairs(int n) {
    if(n == 1 || n == 2) {
        return n;
    } 
    int[] memo = new int[n + 1];
    memo[1] = 1;
    memo[2] = 2;    
    for(int i = 3; i < n + 1; i ++) {
        memo[i] = memo[i-1] + memo[i-2];
    }
    return memo[n];
}
```

### Python

### Go

## 75  Sort Colors

### Java



### Python

### Go