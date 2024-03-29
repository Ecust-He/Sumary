[TOC]

## 1  ~~Two Sum~~

### 思路

- 使用Map数据结构，key为target - nums[i]，value为i

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

- 注意加法进位

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

- 首先需要定义[(eft, right]区间内保存不重复字符串
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

## 4  Median of Two Sorted Arrays

### 思路

- 需要开辟新的空间保存合并之后的数组
- 循环不变式：两个数组至少有一个未遍历结束

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

## 5  Longest Palindromic Substring

### 思路

### Java

### Python

### Go

## 7  ~~Reverse Integer~~

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

## 9  ~~Palindrome Number~~

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

- 指针对撞
- 已短板一边为准

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

## 13  ~~Roman to Integer~~

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
            if (commonPrefix.length() == 0) {
                return "";
            }
            commonPrefix = commonPrefix.substring(0, commonPrefix.length() - 1);
        }
    }
    return commonPrefix;
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

- 先排序
- i，j，k三个数，每个位置的数字不能重复

### Java

```java
public List<List<Integer>> threeSum(int[] nums) {
    int len = nums.length;
    List<List<Integer>> res = new ArrayList<>();
    if(len < 3) {
        return res;
    }
    Arrays.sort(nums);
    for (int i = 0; i < len - 2; i++) {
        if(i > 0 && nums[i] == nums[i-1]) {
            continue;
        }
        int j = i + 1;
        int k = len -1;
        while (j < k) {
            if(nums[i] + nums[j] + nums[k] == 0) {
                res.add(Arrays.asList(nums[i], nums[j], nums[k]));
                j++;
                k--;
                while (j < k && nums[j] == nums[j-1]) {
                    j++;
                }
                while (j < k && nums[k] == nums[k+1]){
                    k--;
                }
            } else if(nums[i] + nums[j] + nums[k] < 0) {
                j++;
            } else {
                k--;
            }
        }
    }
    return res;
}
```

### Python

### Go

## 16  3Sum Closest

### Java

### Python

### Go

## 17  ~~Letter Combinations of a Phone Number~~

### 思路

- 逐一翻译字符串中的每一个字符，并保存

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
        if (index == digits.length()) {
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

- 求出链表的长度
- 找出待删除节点的前一个节点

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

### 思路

循环不变式：至少有一个链表未遍历完成

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
        if(head == null) {
            return null;
        }
        ListNode dumyHead = new ListNode(0, head);
        ListNode pre = dumyHead;
        ListNode node1;
        ListNode node2;
        ListNode next;
        while (pre.next != null && pre.next.next != null) {
            node1 = pre.next;
            node2 = node1.next;
            next = node2.next;

            pre.next = node2;
            node2.next = node1;
            node1.next = next;
            pre = node1;
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

## 26  Remove Duplicates from Sorted Array

### 思路

- 保证[0, len)区间内数字不重复

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

## 46  ~~Permutations~~

#### 思路

- 对于每个元素可选或不选
- 回溯

### Java

```java
public List<List<Integer>> permute(int[] nums) {
    ArrayList<List<Integer>> res = new ArrayList<>();
    permute(nums, new LinkedList<>(), res);
    return res;
}

public void permute(int[] nums, LinkedList<Integer> selected, ArrayList<List<Integer>> res) {
    if(selected.size() == nums.length) {
        res.add(new ArrayList<>(selected));
        return;
    }
    for (int i = 0; i < nums.length; i++) {
        if (!selected.contains(nums[i])) {
            selected.add(nums[i]);
            permute(nums, selected, res);
            selected.removeLast();
        }
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

## 67  ~~Add Binary~~

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
        sum %= 2;
        stack.add(sum);
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

## 70  ~~Climbing Stairs~~

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

## 75  ~~Sort Colors~~

### Java

#### 方法一

```java
public void sortColors(int[] nums) {
    int[] freq = new int[3];
    for (int num : nums) {
        freq[num]++;
    }
    int i = 0;
    for (int j = 0; j < 3; j++) {
        while (freq[j] != 0) {
            nums[i++] = j;
            freq[j]--;
        }
    }
}
```

#### 方法二

```java
public void sortColors(int[] nums) {
    int n = nums.length;
    int left = 0; //[0, left)区间内保存数字0
    int right = n-1;//(right,n-1]区间内保存数字2
    int i = 0;
    while (i <= right) {
        if(nums[i] == 1) {
            i++;
        } else if (nums[i] < 1) {
            swap(nums, i, left);
            left++;
            i++;
        } else {
            swap (nums, i, right);
            right--;
        }
    }
}

public void swap(int[] nums, int i, int j) {
    int temp = nums[i];
    nums[i] = nums[j];
    nums[j] = temp;
}
```

### Python

```python
def sortColors(self, nums):
    """
    :type nums: List[int]
    :rtype: None Do not return anything, modify nums in-place instead.
    """
    n = len(nums)
    left = 0
    right = n - 1
    i = 0
    while i <= right:
        if nums[i] == 1:
            i += 1
            elif nums[i] < 1:
                nums[i], nums[left] = nums[left], nums[i]
                left += 1
                i += 1
                else:
                    nums[i], nums[right] = nums[right], nums[i]
                    right -= 1
```

### Go

```go
func sortColors(nums []int)  {
    n := len(nums)
    i := 0
    left := 0
    right := n-1
    for ;i <= right; {
       if nums[i] == 1 {
          i += 1
      } else {
         if nums[i] == 0 {
            nums[i], nums[left] = nums[left], nums[i]
            left += 1
            i += 1
         } else {
            nums[i], nums[right] = nums[right], nums[i]
            right -= 1
         }
      }
   }
}
```

## 77  Combinations

### Java

```java
public List<List<Integer>> combine(int n, int k) {
    ArrayList<List<Integer>> res = new ArrayList<>();
    combine(n, k, 1, new LinkedList<>(), res);
    return res;
}

/**
 * 从[start, n]区间内选取k个元素保存在selected中
 * @param n
 * @param k
 * @param index
 * @param selected
 * @param res
 */
public void combine(int n, int k, int start, LinkedList<Integer> selected, ArrayList<List<Integer>> res) {
    if(k == selected.size()) {
        res.add(new ArrayList<>(selected));
        return;
    }
    for (int i = start; i <= n; i++) {
        selected.add(i);
        combine(n, k, i + 1, selected, res);
        selected.removeLast();
    }
}
```

### Python

### Go

## 78  Subsets

### Java

### Python

### Go

## 79  Word Search

### Java

```java
public boolean exist(char[][] board, String word) {
    int m = board.length;
    int n = board[0].length;
    visited = new boolean[m][n];

    for (int i = 0; i < m; i++) {
        for (int j = 0; j < n; j++) {
            if (serchWord(board, m, n, i, j, 0, word)) {
                return true;
            }
        }
    }
    return false;
}

/**
 * 考察(i, j)位置是否匹配字符word[index]
 *
 * @param board
 * @param m
 * @param n
 * @param x
 * @param y
 * @param index
 * @param word
 * @param visited
 * @return
 */
public boolean serchWord(char[][] board, int m, int n, int x, int y, int index, String word) {
    if (index == word.length() - 1) {
        return board[x][y] == word.charAt(index);
    }
    if (board[x][y] != word.charAt(index)) {
        return false;
    }
    visited[x][y] = true;
    for (int k = 0; k < 4; k++) {
        int newX = x + d[k][0];
        int newY = y + d[k][1];
        if (inArea(newX, newY, m, n) && !visited[newX][newY] && serchWord(board, m, n, newX, newY, index + 1, word)) {
            return true;
        }
    }
    visited[x][y] = false;
    return false;
}


public boolean inArea(int x, int y, int m, int n) {
    return x >= 0 && x < m && y >= 0 && y < n;
}
```

### Python

### Go

## 82  Remove Duplicates from Sorted List II

### Java

### Python

### Go  

## 83  ~~Remove Duplicates from Sorted List~~

### Java

```java
public ListNode deleteDuplicates(ListNode head) {
    if(head == null || head.next == null) {
        return head;
    }
    ListNode dumyHead = new ListNode(0, head);
    ListNode pre = head;
    ListNode cur = head.next;
    while (cur != null) {
        if(cur.val != pre.val) {
            pre.next = cur;
            pre = cur;
        }
        cur = cur.next;
    }
    pre.next = null;
    return dumyHead.next;
}
```

### Python

### Go

## 88  Merge Sorted Array

### Java



### Python

### Go

## 91  Decode Ways

### Java

### Python

### Go

## 92  Reverse Linked List II

### Java

```java
public ListNode reverseBetween(ListNode head, int left, int right) {
    if(head == null || head.next == null) {
        return head;
    }
    ListNode dumyHead = new ListNode(0, head);
    ListNode leftPre = dumyHead;
    for (int i = 0; i < left -1; i++) {
        leftPre = leftPre.next;
    }
    ListNode pre = leftPre.next;
    ListNode cur = pre.next;
    ListNode leftNode = leftPre.next;
    ListNode next = null;
    for (int i = 0; i < right - left; i++) {
        next = cur.next;
        cur.next = pre;
        pre = cur;
        cur = next;
    }
    leftPre.next = pre;
    leftNode.next = cur;
    return dumyHead.next;
}
```

### Python

### Go

## 94  ~~Binary Tree Inorder Traversal~~

### Java

```java
public List<Integer> inorderTraversal(TreeNode root) {
    if(root == null) {
        return new ArrayList<>();
    }
    ArrayList<Integer> res = new ArrayList<>();
    res.addAll(inorderTraversal(root.left));
    res.add(root.val);
    res.addAll(inorderTraversal(root.right));
    return res;
}
```

### Python

### Go

## 97  Interleaving String

### Java

### Python

### Go

## 100  Same Tree

### Java

### Python

### Go

## 102  ~~Binary Tree Level Order Traversal~~

### Java

```java
public List<List<Integer>> levelOrder(TreeNode root) {
    ArrayList<List<Integer>> res = new ArrayList<>();
    if(root == null) {
        return res;
    }
    Queue<TreeNode> q = new LinkedList<>();
    q.add(root);
    while (!q.isEmpty()) {
        ArrayList<Integer> levelList = new ArrayList<>();
        final int levelSize = q.size();
        for (int i = 0; i < levelSize; i++) {
            TreeNode node = q.poll();
            levelList.add(node.val);
            if(node.left !=null) {
                q.add(node.left);
            }
            if(node.right != null) {
                q.add(node.right);
            }
        }
        res.add(levelList);
    }
    return res;
}
```

### Python

### Go

## 104  ~~Maximum Depth of Binary Tree~~

### Java

```java
public int maxDepth(TreeNode root) {
    if (root == null) {
        return 0;
    }
    return 1 + Math.max(maxDepth(root.left), maxDepth(root.right));
}
```

### Python

### Go

## 111  Minimum Depth of Binary Tree

### Java

### Python

### Go

## 112  ~~Path Sum~~

### Java

```java
public boolean hasPathSum(TreeNode root, int targetSum) {
    if(root == null) {
        return false;
    }
    if(root.left == null && root.right == null) {
        return targetSum == root.val;
    }
    return hasPathSum(root.left, targetSum - root.val) || hasPathSum(root.right, targetSum - root.val);
}
```

### Python

### Go

## 125  ~~Valid Palindrome~~

### Java

```java
public boolean isPalindrome(String s) {
    s = s.toLowerCase();
    int n = s.length();
    int i = 0;
    int j = n - 1;
    while (i < j) {
        while (i < j && !isAlphanumericCharacter(s.charAt(i))) {
            i++;
        }
        while (j > i && !isAlphanumericCharacter(s.charAt(j))) {
            j--;
        }
        if(i >= j) {
            break;
        }
        if(s.charAt(i) != s.charAt(j)){
            return false;
        }
        i++;
        j--;
    }
    return true;
}

public boolean isAlphanumericCharacter(Character c) {
    return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9');
}
```

### Python

### Go

## 132  Palindrome Partitioning II

### Java

### Python

### Go

## 139  Word Break

### Java

### Python

### Go

## 143  Reorder List

### Java

### Python

### Go

## 144  Minimum Depth of Binary Tree

### Java

### Python

### Go

## 145  ~~Binary Tree Postorder Traversal~~

### Java

```java
public List<Integer> postorderTraversal(TreeNode root) {
    ArrayList<Integer> res = new ArrayList<>();
    if(root != null) {
        res.addAll(postorderTraversal(root.left));
        res.addAll(postorderTraversal(root.right));
        res.add(root.val);

    }
    return res;
}
```

### Python

### Go

## 150  Evaluate Reverse Polish Notation

### Java



### Python

### Go

## 160  Intersection of Two Linked Lists

### Java

### Python

### Go

## 167  Two Sum II - Input array is sorted

### Java

```java
public int[] twoSum(int[] numbers, int target) {
    int i = 0;
    int j = numbers.length - 1;
    while(i < j) {
        if(numbers[i] + numbers[j] == target) {
            return new int[]{i+1, j+1};
        } else if(numbers[i] + numbers[j] < target) {
            i++;
        } else {
            j--;
        }
    }
    return null;
}
```

### Python

### Go

## 168  Excel Sheet Column Title

### Java

### Python

### Go

## 175  Combine Two Tables

```mysql
SELECT p.FirstName, p.LastName, a.City, a.State FROM Person p LEFT JOIN Address a
ON p.PersonId = a.PersonId;
```

## 176  Second Highest Salary

```mysql
SELECT max(Salary) As SecondHighestSalary FROM Employee where Salary < (SELECT max(Salary) AS highestSalary FROM Employee);

SELECT max(Salary) As SecondHighestSalary FROM Employee where Salary not in (SELECT max(Salary) AS highestSalary FROM Employee);

SELECT (SELECT DISTINCT Salary FROM Employee ORDER BY Salary DESC LIMIT 1,1) AS SecondHighestSalary;
```

## 177  Nth Highest Salary

```mysql
CREATE FUNCTION getNthHighestSalary(N INT) RETURNS INT
BEGIN
DECLARE M INT;
SET M=N-1;
  RETURN (
      # Write your MySQL query statement below.
      SELECT IFNULL((SELECT DISTINCT Salary FROM Employee ORDER BY Salary DESC LIMIT M ,1), NULL)
  );
END
```

## 178  Rank Scores

```sql
# 通过count函数
select
    Score,(
           select 
                count(distinct(score))
            from 
                Scores
            where
                score >= s1.score
            ) as 'rank'
from
    Scores as s1
order by 
    Score desc
```

## 180  Consecutive Numbers

```mysql
SELECT DISTINCT num AS ConsecutiveNums from Logs AS l1 WHERE num = (SELECT num from Logs where id = l1.id + 1) and num = (SELECT num from Logs where id = l1.id + 2)

SELECT DISTINCT l1.num AS ConsecutiveNums from Logs AS l1, Logs AS l2, Logs AS l3 WHERE l1.id + 1 = l2.id and l1.num = l2.num and l1.id + 2 = l3.id and l1.num = l3.num
```

## 181   Employees Earning More Than Their Managers

```mysql
SELECT E1.name as Employee FROM Employee AS E1, Employee AS E2 WHERE E1.managerId = E2.id AND E1.salary > E2.salary;
```

## 182  Duplicate Emails

```mysql
select distinct p1.email as Email from Person AS p1 where (select count(id) from Person where email = p1.email) >= 2;

select email AS Email from Person group by email having count(*) > 1;
```

## 183  Customers Who Never Order

```mysql
select name AS Customers from Customers where id not in (select o.customerId from Orders AS o left Join Customers AS c ON c.id = o.customerId);
```

## 184  Department Highest Salary

```mysql
select d2.name AS Department, e1.name AS Employee, e1.salary AS Salary from (select d.id, max(e.salary) AS maxSalary from Employee AS e left join Department AS d ON e.departmentId = d.id group by d.id) AS d1, Employee AS e1, Department AS d2 
where d1.id = e1.departmentId and d1.maxSalary = e1.salary and d1.id = d2.id

select d2.name AS Department, e1.name AS Employee, e1.salary AS Salary from Employee AS e1 
inner join (select d.id, max(e.salary) AS maxSalary from Employee AS e left join Department AS d ON e.departmentId = d.id group by d.id) AS d1
inner join Department AS d2 
on d1.id = e1.departmentId and d1.maxSalary = e1.salary and d1.id = d2.id
```

## 185  Department Top Three Salaries

```mysql

```

## 196  Delete Duplicate Emails

```mysql
Delete p1.* from Person as p1,Person as p2 where p1.Email=p2.Email and p1.Id>p2.Id

delete from Person where id not in(select minId from (select min(id) AS minId, email from Person group by email) AS p)
```

## 197  Rising Temperature

```mysql
SELECT w1.Id FROM weather w1,weather w2 WHERE TO_DAYS(w1.Date)=TO_DAYS(w2.Date)+1
AND w1.Temperature>w2.Temperature
```



## 206  ~~Reverse Linked List~~

### Java

#### 方法一

```java
// 递归写法
public ListNode reverseList(ListNode head) {
    if(head == null || head.next == null) {
        return head;
    }
    ListNode newHead = reverseList(head.next);
    head.next.next = head;
    head.next = null;
    return newHead;
}
```

#### 方法二

```java
// 迭代写法
public ListNode reverseList(ListNode head) {
    if(head == null || head.next == null) {
        return head;
    }
    ListNode pre = head;
    ListNode cur = pre.next;
    ListNode next;
    while (cur != null) {
        next = cur.next;
        cur.next = pre;
        pre = cur;
        cur = next;
    }
    head.next = null;
    return pre;
}
```

### Python

### Go

## 217  ~~Contains Duplicate~~

### Java

```java
public boolean containsDuplicate(int[] nums) {
    Set<Integer> set = new HashSet<>();
    for(int num: nums) {
        if(!set.add(num)) {
            return true;
        }
    }
    return false;
}
```

### Python

### Go

## 219  Contains Duplicate II

### Java

```java
public boolean containsNearbyDuplicate(int[] nums, int k) {
    Map<Integer, Integer> map = new HashMap<>();
    for(int i = 0;i < nums.length; i++) {
        int num = nums[i];
        if(map.containsKey(num) && i - map.get(num) <= k) {
            return true;
        }            
        map.put(num, i);
    }
    return false;
}
```

### Python

### Go

## 220  Contains Duplicate III

### Java



### Python

### Go

## 225  Implement Stack using Queues

### Java



### Python

### Go

## 226  Invert Binary Tree

### Java

```java
public TreeNode invertTree(TreeNode root) {
    if(root == null) {
        return null;
    }
    TreeNode left = invertTree(root.right); 
    TreeNode right = invertTree(root.left); 
    root.left = left;
    root.right = right;
    return root;
}
```

### Python

### Go

## 231  Power of Two

### Java

### Python

### Go

## 234  Palindrome Linked List

### Java

### Python

### Go

## 235   Lowest Common Ancestor of a Binary Search Tree

### Java

### Python

### Go

## 237  Delete Node in a Linked List

### Java

### Python

### Go

## 242   ~~Valid Anagram~~

### Java

#### 方法一

```java
public boolean isAnagram(String s, String t) {
    int m = s.length();
    int n = t.length();
    if(m != n) {
        return false;
    }
    Map<Character, Integer> map = new HashMap<>();
    for (int i = 0; i < m; i++) {
        Character c = s.charAt(i);
        map.put(c, map.getOrDefault(c,0) + 1);
    }
    for (int i = 0; i < n; i++) {
        Character c = t.charAt(i);
        if(!map.containsKey(c) || map.get(c) < 1) {
            return false;
        }
        map.put(c, map.get(c) - 1);
    }
    return true;
}
```

#### 方法二

```java
public boolean isAnagram(String s, String t) {
    int[] fre = new int[26];
    int m = s.length();
    int n = t.length();
    if(m != n) {
        return false;
    }
    for (int i = 0; i < m; i++) {
        fre[s.charAt(i) - 'a']++;
        fre[t.charAt(i) - 'a']--;
    }
    for (int i = 0; i < 26; i++) {
        if(fre[i] != 0) {
            return false;
        }
    }
    return true;
}
```

### Python

### Go

## 257  Binary Tree Paths

### Java

```java
public List<String> binaryTreePaths(TreeNode root) {
    List<String> res = new ArrayList<>();
    if(root == null) {
        return res;
    }
    if(root.left == null && root.right == null) {
        res.add(String.valueOf(root.val));
        return res;
    }
    for (String s : binaryTreePaths(root.left)) {
        res.add(root.val + "->" + s);
    }
    for (String s : binaryTreePaths(root.right)) {
        res.add(root.val + "->" + s);
    }
    return res;
}
```

### Python

### Go

## 258  ~~Add Digits~~

### Java

```java
public int addDigits(int num) {
    while (true) {
        int sum = 0;
        while (num > 0) {
            sum += num % 10;
            num /= 10;
        }
        if (sum < 10) {
            return sum;
        }
        num = sum;
    }
}
```

### Python

### Go

## 263  ~~Ugly Number~~

### Java

```java
public boolean isUgly(int n) {
    if(n == 0) {
        return false;
    }
    while (n % 2 == 0) {
        n /= 2;
    }
    while (n % 3 == 0) {
        n /= 3;
    }
    while (n % 5 == 0) {
        n /= 5;
    }
    return n == 1;
}
```

### Python

### Go

## 268  Missing Number

### Java

```java
public int missingNumber(int[] nums) {
    int n = nums.length;
    Set<Integer> set = new HashSet<>();
    for (int num : nums) {
        set.add(num);
    }
    for (int i = 0; i <= n; i++) {
        if(!set.contains(i)){
            return i;
        }
    }
    return -1;
}
```

### Python

### Go

## 279  ~~Perfect Squares~~

### Java

```java
public int numSquares(int n) {
    Queue<Pair> q = new LinkedList<>();
    HashSet<Integer> set = new HashSet<>();
    q.add(new Pair(n, 0));
    while (!q.isEmpty()) {
        Pair p = q.poll();
        int num = p.num;
        int level = p.level;
        for (int i = 0; ; i++) {
            int newNum = num - i * i;
            if(newNum == 0) {
                return level + 1;
            }
            if(newNum < 0) {
                break;
            }
            if(!set.contains(newNum)) {
                q.add(new Pair(newNum,  level + 1));
                set.add(newNum);
            }
        }
    }
    return 0;
}
```

### Python

### Go

## 283  ~~Move Zeroes~~

### Java

```java
public void moveZeroes(int[] nums) {
    int n = nums.length;
    int right = 0;//[0, right)区间内保存非0数字
    for (int i = 0; i < n; i++) {
        if(nums[i] != 0) {
            if(right != i) {
                swap(nums, right, i);
            }
            right++;
        }
    }
}
public void swap(int[] nums, int i, int j) {
    int temp = nums[i];
    nums[i] = nums[j];
    nums[j] = temp;
}
```

### Python

### Go

## 290  Word Pattern

### Java

### Python

### Go

## 300  ~~Longest Increasing Subsequence~~

### Java

```java
public int lengthOfLIS(int[] nums) {
    int n = nums.length;
    if(n == 1) {
        return 1;
    }
    int[] memo = new int[n];// memo[i]记录已nums[i]结尾的最长子序列的长度
    memo[0] = 1;
    for (int i = 1; i < n; i++) {
        memo[i] = 1;
        for (int j = i - 1; j >= 0; j--) {
            if(nums[j] < nums[i]) {
                memo[i] = Math.max(memo[i], 1 + memo[j]);
            }
        }
    }
    int maxLen = 0;
    for (int i = 0; i < n; i++) {
        if(memo[i] > maxLen) {
            maxLen = memo[i];
        }
    }
    return maxLen;
}
```

### Python

### Go

## 326  ~~Power of Three~~

### Java

```java
public boolean isPowerOfThree(int n) {
    if(n == 0) {
        return false;
    }
    while (n % 3 == 0) {
        n /= 3;
    }
    return n == 1;
}
```

### Python

### Go

## 328  Odd Even Linked List

### Java

### Python

### Go

## 337  House Robber III

### Java

### Python

### Go

## 341  Flatten Nested List Iterator

### Java

### Python

### Go

## 343  Integer Break

### Java

### Python

### Go

## 344  ~~Reverse String~~

### Java

```java
public void reverseString(char[] s) {
    int n = s.length;
    int i = 0;
    int j = n - 1;
    while (i < j) {
        swap(s, i, j);
        i++;
        j--;
    }
}

public void swap(char[] s, int i, int j) {
    char c = s[i];
    s[i] = s[j];
    s[j] = c;
}
```

### Python

### Go

## 345  ~~Reverse Vowels of a String~~

### Java

```java
public String reverseVowels(String s) {
    char[] vowels = {'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U'};
    Set<Character> vowelList = new HashSet<>();
    for (int i = 0; i < vowels.length; i++) {
        vowelList.add(vowels[i]);
    }
    char[] chars = s.toCharArray();
    int n = chars.length;
    int i = 0;
    int j = n - 1;
    while (i < j) {
        while (i < j && !vowelList.contains(chars[i])) {
            i++;
        }
        while (i < j && !vowelList.contains(chars[j])) {
            j--;
        }
        if(i >= j) {
            break;
        }
        swap(chars, i, j);
        i++;
        j--;
    }
    StringBuilder sb = new StringBuilder();
    for (int k = 0; k < chars.length; k++) {
        sb.append(chars[k]);
    }
    return sb.toString();
}
public void swap(char[] chars, int i, int j) {
    char c = chars[i];
    chars[i] = chars[j];
    chars[j] = c;
}
```

### Python

### Go

## 347  Top K Frequent Elements

### Java

### Python

### Go

## 349  Intersection of Two Arrays

### Java

### Python

### Go

## 350  Intersection of Two Arrays II

### Java

### Python

### Go

## 376  Wiggle Subsequence

### Java

### Python

### Go

## 392  ~~Is Subsequence~~

### Java

```java
public boolean isSubsequence(String s, String t) {
    int m = s.length();
    int n = t.length();
    int i = 0;
    int j = 0;
    while (i < m && j < n) {
        while (i < m && j < n && t.charAt(j) != s.charAt(i)) {
            j++;
        }
        if(i >= m) {
            return true;
        }
        if(j >= n) {
            break;
        }
        i++;
        j++;
    }
    if(i >= m) {
        return true;
    }
    return false;
}
```

### Python

### Go

## 404  Sum of Left Leaves

### Java

### Python

### Go

## 416  Partition Equal Subset Sum

### Java

### Python

### Go

## 435  Non-overlapping Intervals

### Java

### Python

### Go

## 437  Path Sum III

### Java

### Python

### Go

## 438  Find All Anagrams in a String

### Java

### Python

### Go

## 441  Arranging Coins

### Java

### Python

### Go

## 445  Add Two Numbers II

### Java

### Python

### Go

## 455  ~~Assign Cookies~~

### 思路

贪婪算法

### Java

```java
public int findContentChildren(int[] g, int[] s) {
    int m = g.length;
    int n = s.length;
    int i = 0;
    int j = 0;
    int count = 0;
    Arrays.sort(g);
    Arrays.sort(s);
    while (true) {
        while (i < m && j < n && s[j] < g[i]) {
            j++;
        }
        if (i >= m || j >= n) {
            break;
        }
        count++;
        j++;
        i++;
    }
    return count;
}
```

### Python

### Go

## 474  Ones and Zeroes

### Java

### Python

### Go

## 494  Target Sum

### Java

### Python

### Go

## 507  ~~Perfect Number~~

### Java

```java
public boolean checkPerfectNumber(int num) {
    if(num == 1) {
        return false;
    }
    int sum = 0;
    for (int i = 2; i <= Math.sqrt(num); i++) {
        if(num % i == 0) {
            sum += i + num / i;
        }
    }
    return sum == num - 1;
}
```

### Python

### Go

## 509  Fibonacci Number

### Java

```java
public int fib(int n) {
    if (n == 0 || n == 1) {
        return n;
    }
    int[] memo = new int[n + 1];
    memo[0] = 0;
    memo[1] = 1;
    for (int i = 2; i <= n; i++) {
        memo[i] = memo[i - 1] + memo[i - 2];
    }
    return memo[n];
}
```

### Python

### Go

## 530  ~~Minimum Absolute Difference in BST~~

### Java

```java
public int getMinimumDifference(TreeNode root) {
    List<Integer> res = inOrder(root);
    int min = Integer.MAX_VALUE;
    for (int i = 1; i < res.size(); i++) {
        int diff = res.get(i) - res.get(i-1);
        if(diff < min) {
            min = diff;
        }
    }
    return min;
}

public List<Integer> inOrder(TreeNode root) {
    ArrayList<Integer> res = new ArrayList<>();
    if(root == null) {
        return res;
    }
    res.addAll(inOrder(root.left));
    res.add(root.val);
    res.addAll(inOrder(root.right));
    return res;
}
```

### Python

### Go

## 538  Convert BST to Greater Tree

### Java

### Python

### Go

## 557  Reverse Words in a String III

### Java

```java
public String reverseWords(String s) {
    int n = s.length();
    String[] words = s.split(" ");
    StringBuilder sb = new StringBuilder();
    for (String word : words) {
        int m = word.length();
        for (int i = m - 1; i >= 0; i--) {
            sb.append(String.valueOf(word.charAt(i)));
        }
        sb.append(" ");
    }
    sb.deleteCharAt(sb.length() -1);
    return sb.toString();
}
```

### Python

### Go

## 559  Maximum Depth of N-ary Tree

### Java

### Python

### Go

## 561  Array Partition I

### Java

### Python

### Go

## 595  Big Countries

```mysql

```

## 596  ~~Classes More Than 5 Students~~

```mysql
select class from Courses group by class having count(*) >= 5;
```

## 601  Human Traffic of Stadium

```mysql

```

## 617  Merge Two Binary Trees

### Java

### Python

### Go

## 620  Not Boring Movies

```mysql
select * from Cinema where id %2 = 1 and description != 'boring' order by rating desc;
```

## 627  Swap Salary

```mysql
UPDATE salary 
SET sex = (
    CASE sex
    WHEN 'm' THEN 'f'
    ELSE 'm'
    END
);
```



## 637  ~~Average of Levels in Binary Tree~~

### Java

```java
    ArrayList<Double> res = new ArrayList<>();
    Queue<TreeNode> q = new LinkedList<>();
    q.add(root);
    while (!q.isEmpty()) {
        int size = q.size();
        double levelTotal = 0;
        for (int i = 0; i < size; i++) {
            TreeNode node = q.poll();
            levelTotal += node.val;
            if(node.left != null) {
                q.add(node.left);
            }
            if(node.right != null) {
                q.add(node.right);
            }
        }
        res.add(levelTotal / size);
    }
    return res;
}
```

### Python

### Go

## 645  Set Mismatch

### Java

### Python

### Go

## 653  Two Sum IV - Input is a BST

### Java

```java
public boolean findTarget(TreeNode root, int k) {
    ArrayList<Integer> list = inOrderTraversal(root);
    int n = list.size();
    int i = 0;
    int j = n - 1;
    while (i < j) {
        if(list.get(i) + list.get(j) == k) {
            return true;
        } else if(list.get(i) + list.get(j) < k) {
            i++;
        }else {
            j--;
        }
    }
    return false;
}

public ArrayList<Integer> inOrderTraversal(TreeNode root) {
    ArrayList<Integer> res = new ArrayList<>();
    if(root == null) {
        return res;
    }
    res.addAll(inOrderTraversal(root.left));
    res.add(root.val);
    res.addAll(inOrderTraversal(root.right));
    return res;
}
```

### Python

### Go

## 657  Robot Return to Origin

### Java

```java
public boolean judgeCircle(String moves) {
    int x = 0;
    int y = 0;
    for (int i = 0; i < moves.length(); i++) {
        char c = moves.charAt(i);
        if(c == 'R') {
            x++;
        } else if(c == 'L') {
            x--;
        } else if(c == 'U') {
            y++;
        } else {
            y--;
        }
    }
    return x == 0 && y == 0;
}
```

### Python

### Go

## 677  Map Sum Pairs

### Java

### Python

### Go

## 680  ~~Valid Palindrome II~~

### Java

```java
public boolean validPalindrome(String s) {
    int i = 0;
    int j = s.length() - 1;
    while (i < j) {
        if(s.charAt(i) != s.charAt(j)) {
            break;
        }
        i++;
        j--;
    }
    if(i >=j) {
        return true;
    }
    return isPalindrome(s, i+1, j) || isPalindrome(s, i, j-1);
}

public boolean isPalindrome(String s, int i, int j) {
    while (i < j) {
        if(s.charAt(i) != s.charAt(j)) {
            return false;
        }
        i++;
        j--;
    }
    return true;
}
```

### Python

### Go

## 686  ~~Repeated String Match~~

### Java

```java
public int repeatedStringMatch(String a, String b) {
    int m = a.length();
    int n = b.length();
    int times = n % m == 0 ? n / m : n / m + 1;
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < times; i++) {
        sb.append(a);
    }
    if(sb.toString().indexOf(b) > -1) {
        return times;
    }
    sb.append(a);
    if(sb.toString().indexOf(b) > -1) {
        return times + 1;
    }
    return -1;
}
```

### Python

### Go

## 692  ~~Top K Frequent Words~~

### Java

```java
public List<String> topKFrequent(String[] words, int k) {
    HashMap<String, Integer> fre = new HashMap<>();
    for (String word : words) {
        fre.put(word , fre.getOrDefault(word, 0) + 1);
    }
    Queue<String> q = new PriorityQueue<>((s1, s2) -> {
        if(fre.get(s1) == fre.get(s2)) {
            return s2.compareTo(s1);
        }
        return fre.get(s1) - fre.get(s2);
    });
    for (String word : fre.keySet()) {
        q.add(word);
        if(q.size() > k) {
            q.poll();
        }
    }
    List<String> ret = new ArrayList<>();
    Stack<String> stack = new Stack<>();
    while (!q.isEmpty()) {
        stack.add(q.poll());
    }
    while (!stack.isEmpty()) {
        ret.add(stack.pop());
    }
    return ret;
}
```

### Python

### Go

## 700  ~~Search in a Binary Search Tree~~

### Java

```java
public TreeNode searchBST(TreeNode root, int val) {
    if (root == null) {
        return null;
    }
    if(root.val == val) {
        return root;
    } else if(root.val > val) {
        return searchBST(root.left, val);
    } else {
        return searchBST(root.right, val);
    }
}
```

### Python

### Go

## 703  ~~Kth Largest Element in a Stream~~

### Java

```java
class KthLargest {

    Queue<Integer> q = new PriorityQueue<>();
    int k;
    public KthLargest(int k, int[] nums) {
        this.k = k;
        for (int num : nums) {
            q.add(num);
            if(q.size() > k) {
                q.poll();
            }
        }
    }
    
    public int add(int val) {
        q.add(val);
        if(q.size() > k) {
            q.poll();
        }
        return q.peek();
    }
}
```

### Python

### Go

## 705  Design HashSet

### Java

### Python

### Go

## 709   To Lower Case

### Java

### Python

### Go

## 728  ~~Self Dividing Numbers~~

### Java

```java
public List<Integer> selfDividingNumbers(int left, int right) {
    ArrayList<Integer> res = new ArrayList<>();
    for (int i = left; i <= right; i++) {
        if(isSelfDividingNumber(i)) {
            res.add(i);
        }
    }
    return res;
}

public boolean isSelfDividingNumber(int i) {
    final int num = i;
    while (i > 0) {
        int remain = i % 10;
        if(remain == 0) {
            return false;
        }
        if(num % remain != 0) {
            return false;
        }
        i /= 10;
    }
    return true;
}
```

### Python

### Go

## 771  ~~Jewels and Stones~~

### Java

```java
public int numJewelsInStones(String jewels, String stones) {
    HashSet<Character> set = new HashSet<>();
    for (int i = 0; i < jewels.length(); i++) {
        set.add(jewels.charAt(i));
    }
    int count = 0;
    for (int i = 0; i < stones.length(); i++) {
        if(set.contains(stones.charAt(i))) {
            count++;
        }
    }
    return count;
}
```

### Python

### Go

## 804  ~~Unique Morse Code Words~~

### Java

### Python

### Go

## 832  Flipping an Image

### Java

### Python

### Go

## 844  Backspace String Compare

### Java

### Python

### Go

## 872  ~~Leaf-Similar Trees~~

### Java

```java
public boolean leafSimilar(TreeNode root1, TreeNode root2) {
    List<Integer> list1 = findAllLeaf(root1);
    List<Integer> list2 = findAllLeaf(root2);
    int m = list1.size();
    int n = list2.size();
    if(m != n) {
        return false;
    }
    for (int i = 0; i < m; i++) {
        if(list1.get(i) != list2.get(i)) {
            return false;
        }
    }
    return true;
}

public List<Integer> findAllLeaf(TreeNode root) {
    List<Integer> res = new ArrayList<>();
    if(root == null) {
        return res;
    }
    res.addAll(findAllLeaf(root.left));
    if(root.left == null && root.right == null) {
        res.add(root.val);
        return res;
    }
    res.addAll(findAllLeaf(root.right));
    return res;
}
```

### Python

### Go

## 876  ~~Middle of the Linked List~~

### Java

```java
public ListNode middleNode(ListNode head) {
    ListNode p = head;
    int len = 0;
    while (p != null) {
        len++;
        p = p.next;
    }
    p = head;
    for (int i = 0; i < len / 2; i++) {
        p = p.next;
    }
    return p;
}
```

### Python

### Go

## 897  ~~Increasing Order Search Tree~~

### Java

```java
public TreeNode increasingBST(TreeNode root) {
    List<Integer> res = inOrderTraversal(root);
    TreeNode dumyHead = new TreeNode(0);
    TreeNode p = dumyHead;
    for (Integer num : res) {
        p.right = new TreeNode(num);
        p = p.right;
    }
    return dumyHead.right;
}

public List<Integer> inOrderTraversal(TreeNode root) {
    List<Integer> res = new ArrayList<>();
    if(root == null) {
        return res;
    }
    res.addAll(inOrderTraversal(root.left));
    res.add(root.val);
    res.addAll(inOrderTraversal(root.right));
    return res;
}
```

### Python

### Go

## 905  ~~Sort Array By Parity~~

### Java

```java
public int[] sortArrayByParity(int[] nums) {
    int n = nums.length;
    int left = n - 1;//(left,n)区间内所有数都是奇数
    int i = 0;
    while (true) {
        while (i < left && nums[i] % 2 == 0) {
            i++;
        }
        if(i >= left) {
            return nums;
        }
        swap(nums, i, left);
        left--;
    }

}

public void  swap(int[] nums, int i, int j) {
    int t = nums[i];
    nums[i] = nums[j];
    nums[j] = t;
}
```

### Python

### Go

## 922  Sort Array By Parity II

### Java

### Python

### Go

## 938  ~~Range Sum of BST~~

### Java

#### 方法一

```java
public int rangeSumBST(TreeNode root, int low, int high) {
    int sum = 0;
    if(root == null) {
        return 0;
    }
    Stack<TreeNode> stack = new Stack<>();
    stack.add(root);
    int val = root.val;
    if(val >= low && val <= high) {
        sum += val;
    }
    while (!stack.isEmpty()) {
        TreeNode p = stack.pop();
        if(p.left != null) {
            stack.add(p.left);
            val = p.left.val;
            if(val >= low && val <= high) {
                sum += val;
            }
        }
        if(p.right != null) {
            stack.add(p.right);
            val = p.right.val;
            if(val >= low && val <= high) {
                sum += val;
            }
        }
    }
    return sum;
}
```

#### 方法二

```java
public int rangeSumBST(TreeNode root, int low, int high) {
    if(root == null) {
        return 0;
    }
    int sum = 0;
    sum += rangeSumBST(root.left, low, high);
    if(root.val >= low &&root.val <= high) {
        sum += root.val;
    }
    sum += rangeSumBST(root.right, low, high);
    return sum;
}
```

### Python

### Go

## 965  ~~Univalued Binary Tree~~

### Java

```java
public boolean isUnivalTree(TreeNode root) {
    if(root == null) {
        return true;
    }
    if(root.left != null) {
        if(root.val != root.left.val) {
            return false;
        }
    }
    if(root.right != null) {
        if(root.val != root.right.val) {
            return false;
        }
    }
    return isUnivalTree(root.left) && isUnivalTree(root.right);
}
```

### Python

### Go

## 977  Squares of a Sorted Array

### Java

```java
public int[] sortedSquares(int[] nums) {
    int n = nums.length;
    int[] ret = new int[n];
    for (int i = 0; i < n; i++) {
        int num = nums[i];
        ret[i] = num * num;
    }
    Arrays.sort(ret);
    return ret;
}
```

### Python

### Go

## 1002  Find Common Characters

### Java

### Python

### Go

## 1019  Next Greater Node In Linked List

### Java

### Python

### Go

## 1022   Sum of Root To Leaf Binary Numbers

### Java

### Python

### Go

## 1114  print in Order

### Java

### Python

### Go

## 1122  Relative Sort Array

### Java

### Python

### Go

## 1160  Find Words That Can Be Formed by Characters

### Java

### Python

### Go

## 1200  Minimum Absolute Difference

### Java

### Python

### Go

## 1221  Split a String in Balanced Strings

### Java

```java
public int balancedStringSplit(String s) {
    int count = 0;
    int RNum = 0;
    int LNum = 0;
    for (int i = 0; i < s.length(); i++) {
        Character c = s.charAt(i);
        if(RNum == LNum) {
            count ++;
        }
        if(c == 'R') {
            RNum++;
        }
        if(c == 'L') {
            LNum++;
        }
    }
    return count;
}
```

### Python

### Go

## 1290  ~~Convert Binary Number in a Linked List to Integer~~

### Java

```java
public int getDecimalValue(ListNode head) {
    Stack<Integer> stack = new Stack<>();
    ListNode p = head;
    while (p != null) {
        stack.add(p.val);
        p = p.next;
    }
    int sum = 0;
    int i = 0;
    while (!stack.isEmpty()) {
        sum += stack.pop() * Math.pow(2, i);
        i++;
    }
    return sum;
}
```

### Python

### Go

## 1295  ~~Find Numbers with Even Number of Digits~~

### Java



### Python

### Go

## 1299  Replace Elements with Greatest Element on Right Side

### Java

### Python

### Go

## 1304  Find N Unique Integers Sum up to Zero

### Java

### Python

### Go

## 1309  Decrypt String from Alphabet to Integer Mapping

### Java

### Python

### Go

## 1313  Decompress Run-Length Encoded List

### Java

### Python

### Go

## 1337  The K Weakest Rows in a Matrix

### Java

### Python

### Go

## 1351  Count Negative Numbers in a Sorted Matrix

### Java

### Python

### Go

## 1365  How Many Numbers Are Smaller Than the Current Number

### Java

```java
public int[] smallerNumbersThanCurrent(int[] nums) {
    int n = nums.length;
    int[] ret = new int[n];
    for (int i = 0; i < n; i++) {
        int count = 0;
        for (int j = 0; j < n; j++) {
            if(nums[j] < nums[i] && i != j) {
                count++;
            }
        }
        ret[i] = count;
    }
    return ret;
}
```

### Python

### Go

## 1385  Find the Distance Value Between Two Arrays

### Java

### Python

### Go

## 1394  Find Lucky Integer in an Array

### Java

```java
public int findLucky(int[] arr) {
    HashMap<Integer, Integer> fre = new HashMap<>();
    for (int num : arr) {
        fre.put(num, fre.getOrDefault(num, 0) + 1);
    }
    int luckyNum = -1;
    for (Map.Entry<Integer, Integer> entry : fre.entrySet()) {
        if(entry.getKey() == entry.getValue()) {
            if(entry.getKey() > luckyNum) {
                luckyNum = entry.getKey();
            }
        }
    }
    return luckyNum;
}
```

### Python

### Go

## 1413  Minimum Value to Get Positive Step by Step Sum

### Java

### Python

### Go

## 1436  Destination City

### Java

### Python

### Go
