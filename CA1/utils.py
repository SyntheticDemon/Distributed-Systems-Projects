def read_orders_from_file(file_path:str):
    names = []
    try:
        with open(file_path, 'r') as file:
            for line in file:
                # Assuming each line contains a single name
                name = line.strip()
                names.append(name)
        return names
    except FileNotFoundError:
        print(f"File '{file_path}' not found.")
        return []
    
    
    
def find_orders_with_prefix(prefix, orders):
    return [order for order in orders if order.startswith(prefix)]