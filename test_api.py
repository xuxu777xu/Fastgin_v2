import requests
import json
from colorama import init, Fore, Style
import time

# 初始化 colorama
init()

# 服务器配置
BASE_URL = "http://localhost:8080/api"

def print_response(response, test_name):
    """打印响应结果"""
    print(f"\n{Fore.CYAN}=== {test_name} ==={Style.RESET_ALL}")
    print(f"{Fore.YELLOW}状态码:{Style.RESET_ALL}", response.status_code)
    print(f"{Fore.YELLOW}响应头:{Style.RESET_ALL}", dict(response.headers))
    print(f"{Fore.YELLOW}响应体:{Style.RESET_ALL}", json.dumps(response.json(), ensure_ascii=False, indent=2))

def test_get_user():
    """测试获取用户接口"""
    # 测试正确的请求
    response = requests.get(f"{BASE_URL}/user", params={"id": "1"})
    print_response(response, "测试获取用户 - 正确请求")

    # 测试缺少参数的请求
    response = requests.get(f"{BASE_URL}/user")
    print_response(response, "测试获取用户 - 缺少参数")

def test_create_user():
    """测试创建用户接口"""
    # 测试正确的请求
    data = {
        "username": "李四",
        "age": 30
    }
    response = requests.post(f"{BASE_URL}/user", json=data)
    print_response(response, "测试创建用户 - 正确请求")

    # 测试缺少必要字段的请求
    data = {
        "username": "李四"
        # 缺少 age 字段
    }
    response = requests.post(f"{BASE_URL}/user", json=data)
    print_response(response, "测试创建用户 - 缺少年龄字段")

    # 测试年龄超出范围的请求
    data = {
        "username": "李四",
        "age": 200  # 超出范围
    }
    response = requests.post(f"{BASE_URL}/user", json=data)
    print_response(response, "测试创建用户 - 年龄超出范围")

def test_rate_limit():
    """测试限流功能"""
    print(f"\n{Fore.CYAN}=== 测试限流功能 ==={Style.RESET_ALL}")
    start_time = time.time()
    success_count = 0
    fail_count = 0
    
    # 快速发送300个请求
    for i in range(300):
        response = requests.get(f"{BASE_URL}/user", params={"id": "1"})
        if response.status_code == 200:
            success_count += 1
        else:
            fail_count += 1
    
    end_time = time.time()
    duration = end_time - start_time
    
    print(f"{Fore.GREEN}成功请求:{Style.RESET_ALL}", success_count)
    print(f"{Fore.RED}失败请求:{Style.RESET_ALL}", fail_count)
    print(f"{Fore.YELLOW}总耗时:{Style.RESET_ALL}", f"{duration:.2f}秒")

def test_login():
    """测试登录接口"""
    url = f"{BASE_URL}/login"
    data = {
        "username": "testuser",
        "password": "testpass"
    }
    response = requests.post(url, json=data)
    print("\n=== 测试登录接口 ===")
    print(f"状态码: {response.status_code}")
    print(f"响应内容: {response.json()}")

def test_unpaid_orders():
    """测试获取未支付订单接口"""
    url = f"{BASE_URL}/unpaid-orders"
    data = {
        "user_id": "12345",
        "page": 1,
        "size": 10
    }
    response = requests.post(url, json=data)
    print("\n=== 测试获取未支付订单接口 ===")
    print(f"状态码: {response.status_code}")
    print(f"响应内容: {response.json()}")

def test_payment():
    """测试支付接口"""
    url = f"{BASE_URL}/payment"
    data = {
        "order_id": "ORDER123",
        "amount": 99.99,
        "payment_type": "alipay"
    }
    response = requests.post(url, json=data)
    print("\n=== 测试支付接口 ===")
    print(f"状态码: {response.status_code}")
    print(f"响应内容: {response.json()}")

def test_hmf_ci():
    """测试获取hmfCi接口"""
    url = f"{BASE_URL}/hmf-ci"
    data = {
        "user_id": "12345",
        "token": "test-token"
    }
    response = requests.post(url, json=data)
    print("\n=== 测试获取hmfCi接口 ===")
    print(f"状态码: {response.status_code}")
    print(f"响应内容: {response.json()}")

def run_all_tests():
    """运行所有测试"""
    try:
        print(f"{Fore.GREEN}开始测试 API{Style.RESET_ALL}")
        
        # 测试获取用户接口
        test_get_user()
        
        # 测试创建用户接口
        test_create_user()
        
        # 测试限流功能
        test_rate_limit()
        
        test_login()
        test_unpaid_orders()
        test_payment()
        test_hmf_ci()
        
        print(f"\n{Fore.GREEN}测试完成{Style.RESET_ALL}")
        
    except requests.exceptions.ConnectionError:
        print(f"{Fore.RED}错误: 无法连接到服务器，请确保服务器已启动并运行在正确的端口上{Style.RESET_ALL}")
    except Exception as e:
        print(f"{Fore.RED}错误: {str(e)}{Style.RESET_ALL}")

if __name__ == "__main__":
    print("开始测试 API 接口...")
    run_all_tests()
    print("\n测试完成!") 