import requests
import json
from colorama import init, Fore, Style, Back
import time
import random
from concurrent.futures import ThreadPoolExecutor
import sys

# 初始化 colorama
init()

# 服务器配置
BASE_URL = "http://localhost:8080/api"

class APITester:
    def __init__(self):
        self.session = requests.Session()
        self.success_count = 0
        self.fail_count = 0
        self.test_results = []

    def print_test_header(self, test_name):
        """打印测试标题"""
        print(f"\n{Back.BLUE}{Fore.WHITE} 测试: {test_name} {Style.RESET_ALL}")

    def print_response(self, response, test_case):
        """打印响应结果"""
        is_success = 200 <= response.status_code < 300
        status_color = Fore.GREEN if is_success else Fore.RED
        
        print(f"\n{Fore.CYAN}【测试用例】{test_case}{Style.RESET_ALL}")
        print(f"{Fore.YELLOW}状态码:{Style.RESET_ALL} {status_color}{response.status_code}{Style.RESET_ALL}")
        
        try:
            response_json = response.json()
            print(f"{Fore.YELLOW}响应体:{Style.RESET_ALL}\n{json.dumps(response_json, ensure_ascii=False, indent=2)}")
        except json.JSONDecodeError:
            print(f"{Fore.RED}响应体不是有效的JSON格式{Style.RESET_ALL}")
            print(f"{Fore.YELLOW}原始响应:{Style.RESET_ALL}\n{response.text}")

        if is_success:
            self.success_count += 1
        else:
            self.fail_count += 1

        # 记录测试结果
        self.test_results.append({
            "test_case": test_case,
            "status_code": response.status_code,
            "is_success": is_success,
            "response": response_json if 'response_json' in locals() else response.text
        })

    def test_example_apis(self):
        """测试示例API接口"""
        self.print_test_header("示例API测试")

        # 测试获取用户
        try:
            # 正常请求
            response = self.session.get(f"{BASE_URL}/user", params={"id": "1"})
            self.print_response(response, "获取用户 - 正常请求")

            # 缺少ID参数
            response = self.session.get(f"{BASE_URL}/user")
            self.print_response(response, "获取用户 - 缺少ID参数")

            # ID参数无效
            response = self.session.get(f"{BASE_URL}/user", params={"id": "invalid"})
            self.print_response(response, "获取用户 - 无效ID参数")
        except Exception as e:
            print(f"{Fore.RED}测试获取用户接口时发生错误: {str(e)}{Style.RESET_ALL}")

        # 测试创建用户
        try:
            # 正常请求
            data = {
                "username": "测试用户",
                "age": 25
            }
            response = self.session.post(f"{BASE_URL}/user", json=data)
            self.print_response(response, "创建用户 - 正常请求")

            # 缺少必要字段
            data = {
                "username": "测试用户"
            }
            response = self.session.post(f"{BASE_URL}/user", json=data)
            self.print_response(response, "创建用户 - 缺少年龄字段")

            # 无效的年龄值
            data = {
                "username": "测试用户",
                "age": -1
            }
            response = self.session.post(f"{BASE_URL}/user", json=data)
            self.print_response(response, "创建用户 - 无效年龄值")
        except Exception as e:
            print(f"{Fore.RED}测试创建用户接口时发生错误: {str(e)}{Style.RESET_ALL}")

    def test_login_api(self):
        """测试登录接口"""
        self.print_test_header("登录接口测试")

        try:
            # 正常登录（带自定义ID）
            data = {
                "username": "testuser",
                "password": "testpass"
            }
            headers = {
                'X-Request-ID': 'test-request-id',
                'X-Trace-ID': 'test-trace-id'
            }
            response = self.session.post(f"{BASE_URL}/login", json=data, headers=headers)
            self.print_response(response, "登录 - 正常请求（带自定义ID）")
            
            # 验证响应头中的ID
            print(f"{Fore.YELLOW}响应头中的请求ID:{Style.RESET_ALL}", response.headers.get('X-Request-ID'))
            print(f"{Fore.YELLOW}响应头中的追踪ID:{Style.RESET_ALL}", response.headers.get('X-Trace-ID'))

            # 正常登录（自动生成ID）
            response = self.session.post(f"{BASE_URL}/login", json=data)
            self.print_response(response, "登录 - 正常请求（自动生成ID）")
            
            # 验证响应头中是否包含自动生成的ID
            print(f"{Fore.YELLOW}自动生成的请求ID:{Style.RESET_ALL}", response.headers.get('X-Request-ID'))
            print(f"{Fore.YELLOW}自动生成的追踪ID:{Style.RESET_ALL}", response.headers.get('X-Trace-ID'))

            # 缺少密码
            data = {
                "username": "testuser"
            }
            response = self.session.post(f"{BASE_URL}/login", json=data)
            self.print_response(response, "登录 - 缺少密码")

            # 缺少用户名
            data = {
                "password": "testpass"
            }
            response = self.session.post(f"{BASE_URL}/login", json=data)
            self.print_response(response, "登录 - 缺少用户名")

            # 空的请求体
            response = self.session.post(f"{BASE_URL}/login", json={})
            self.print_response(response, "登录 - 空请求体")
        except Exception as e:
            print(f"{Fore.RED}测试登录接口时发生错误: {str(e)}{Style.RESET_ALL}")

    def test_unpaid_orders_api(self):
        """测试获取未支付订单接口"""
        self.print_test_header("未支付订单接口测试")

        try:
            # 正常请求
            data = {
                "user_id": "12345",
                "page": 1,
                "size": 10
            }
            response = self.session.post(f"{BASE_URL}/unpaid-orders", json=data)
            self.print_response(response, "获取未支付订单 - 正常请求")

            # 无效的分页参数
            data = {
                "user_id": "12345",
                "page": 0,
                "size": 0
            }
            response = self.session.post(f"{BASE_URL}/unpaid-orders", json=data)
            self.print_response(response, "获取未支付订单 - 无效分页参数")

            # 缺少用户ID
            data = {
                "page": 1,
                "size": 10
            }
            response = self.session.post(f"{BASE_URL}/unpaid-orders", json=data)
            self.print_response(response, "获取未支付订单 - 缺少用户ID")
        except Exception as e:
            print(f"{Fore.RED}测试获取未支付订单接口时发生错误: {str(e)}{Style.RESET_ALL}")

    def test_payment_api(self):
        """测试支付接口"""
        self.print_test_header("支付接口测试")

        try:
            # 正常请求
            data = {
                "order_id": "ORDER123",
                "amount": 99.99,
                "payment_type": "alipay"
            }
            response = self.session.post(f"{BASE_URL}/payment", json=data)
            self.print_response(response, "支付 - 正常请求")

            # 无效的金额
            data = {
                "order_id": "ORDER123",
                "amount": -1,
                "payment_type": "alipay"
            }
            response = self.session.post(f"{BASE_URL}/payment", json=data)
            self.print_response(response, "支付 - 无效金额")

            # 无效的支付类型
            data = {
                "order_id": "ORDER123",
                "amount": 99.99,
                "payment_type": "invalid_type"
            }
            response = self.session.post(f"{BASE_URL}/payment", json=data)
            self.print_response(response, "支付 - 无效支付类型")
        except Exception as e:
            print(f"{Fore.RED}测试支付接口时发生错误: {str(e)}{Style.RESET_ALL}")

    def test_hmf_ci_api(self):
        """测试获取hmfCi接口"""
        self.print_test_header("hmfCi接口测试")

        try:
            # 正常请求
            data = {
                "user_id": "12345",
                "token": "valid_token"
            }
            response = self.session.post(f"{BASE_URL}/hmf-ci", json=data)
            self.print_response(response, "获取hmfCi - 正常请求")

            # 无效的token
            data = {
                "user_id": "12345",
                "token": ""
            }
            response = self.session.post(f"{BASE_URL}/hmf-ci", json=data)
            self.print_response(response, "获取hmfCi - 无效token")

            # 缺少用户ID
            data = {
                "token": "valid_token"
            }
            response = self.session.post(f"{BASE_URL}/hmf-ci", json=data)
            self.print_response(response, "获取hmfCi - 缺少用户ID")
        except Exception as e:
            print(f"{Fore.RED}测试获取hmfCi接口时发生错误: {str(e)}{Style.RESET_ALL}")

    def test_rate_limit(self):
        """测试限流功能"""
        self.print_test_header("限流功能测试")

        def make_request():
            try:
                response = self.session.get(f"{BASE_URL}/user", params={"id": "1"})
                return response.status_code
            except:
                return 0

        print(f"{Fore.YELLOW}开始并发测试限流功能...{Style.RESET_ALL}")
        start_time = time.time()
        success_count = 0
        fail_count = 0

        # 使用线程池模拟并发请求
        with ThreadPoolExecutor(max_workers=10) as executor:
            futures = [executor.submit(make_request) for _ in range(400)]
            for future in futures:
                status_code = future.result()
                if 200 <= status_code < 300:
                    success_count += 1
                else:
                    fail_count += 1

        duration = time.time() - start_time
        print(f"\n{Fore.GREEN}成功请求: {success_count}{Style.RESET_ALL}")
        print(f"{Fore.RED}失败请求: {fail_count}{Style.RESET_ALL}")
        print(f"{Fore.YELLOW}总耗时: {duration:.2f}秒{Style.RESET_ALL}")

    def print_summary(self):
        """打印测试总结"""
        print(f"\n{Back.WHITE}{Fore.BLACK} 测试总结 {Style.RESET_ALL}")
        print(f"{Fore.GREEN}成功测试用例: {self.success_count}{Style.RESET_ALL}")
        print(f"{Fore.RED}失败测试用例: {self.fail_count}{Style.RESET_ALL}")
        success_rate = (self.success_count / (self.success_count + self.fail_count)) * 100
        print(f"{Fore.YELLOW}成功率: {success_rate:.2f}%{Style.RESET_ALL}")

        # 保存测试结果到文件
        result_file = "test_results.json"
        with open(result_file, "w", encoding="utf-8") as f:
            json.dump({
                "timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
                "success_count": self.success_count,
                "fail_count": self.fail_count,
                "success_rate": success_rate,
                "test_results": self.test_results
            }, f, ensure_ascii=False, indent=2)
        print(f"\n{Fore.CYAN}测试结果已保存到: {result_file}{Style.RESET_ALL}")

    def run_all_tests(self):
        """运行所有测试"""
        try:
            print(f"\n{Back.GREEN}{Fore.BLACK} 开始API测试 {Style.RESET_ALL}")
            
            # 测试示例API
            self.test_example_apis()
            
            # 测试登录相关API
            self.test_login_api()
            self.test_unpaid_orders_api()
            self.test_payment_api()
            self.test_hmf_ci_api()
            
            # 测试限流功能
            self.test_rate_limit()
            
            # 打印测试总结
            self.print_summary()
            
        except requests.exceptions.ConnectionError:
            print(f"\n{Back.RED}{Fore.WHITE} 错误: 无法连接到服务器，请确保服务器已启动并运行在 http://localhost:8080 {Style.RESET_ALL}")
            sys.exit(1)
        except Exception as e:
            print(f"\n{Back.RED}{Fore.WHITE} 测试过程中发生错误: {str(e)} {Style.RESET_ALL}")
            sys.exit(1)

if __name__ == "__main__":
    # 创建测试器实例并运行测试
    tester = APITester()
    tester.run_all_tests() 