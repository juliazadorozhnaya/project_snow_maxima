class dot:
    def __init__(self, x, y, z=0):
        self.x = x
        self.y = y
        self.z = z

    def __eq__(self, d):
        return (self.x==d.x)and(self.y==d.y)

    def __hash__(self):
        return hash((self.x, self.y))
    
    def get_distance2XY(self, d):
        return pow((self.x-d.x), 2) + pow((self.y-d.y), 2)# + pow((self.z-d.z), 2)/18

    def __str__(self):
        return str((self.x, self.y, self.z))

    def __repr__(self):
        return str((self.x, self.y, self.z))

class vect:
    def __init__(self, b, e):
        self.x = b.x - e.x
        self.y = b.y - e.y
        self.z = b.z - e.z

    #не учитывает z координату!
    def __eq__(self, d):
        return (self.x==d.x)and(self.y==d.y)

    def __hash__(self):
        return hash((self.x, self.y))
    
    def len2(self):
        return self.x*self.x + self.y*self.y

    def sprod(self, v):
        return self.x*v.x + self.y*v.y

    def vprod(self, v):
        return self.x*v.y-self.y*v.x

    def add(self, v):
        return vect(self.x+v.x, self.y+v.y)

    def sub(self, v):
        return vect(self.x-v.x, self.y-v.y)

    def mul(self, k):
        return vect(self.x*k, self.y*k)

    def __repr__(self):
        return str((self.x, self.y))

class edge:
    def __init__(self, p1, p2):
        if p1 == p2:
            print("ERROR! IT'S THE SAME POINTS", p1, p2)
            self.e = (p1, p2)
        #возможно плохое решение: вектор ориентирован,а ребро нет. разные смыслы
        self.v = vect(p1, p2)
        #обдумать! может поменять __eq__ вместо упорядочивания
        if p1.x < p2.x:
            self.e = (p1, p2)
        if p1.x > p2.x:
            self.e = (p2, p1)
        if p1.x == p2.x:
            if p1.y < p2.y:
                self.e = (p1, p2)
            if p1.y > p2.y:
                self.e = (p2, p1)
    
    def __eq__(self, d):
        return (self.e==d.e)

    def __hash__(self):
        return hash(self.e)

    def L2norm(self, p):
        return self.e[0].get_distance2XY(p) + self.e[1].get_distance2XY(p)

    def L2len(self,):
        return self.e[0].get_distance2XY(self.e[1])

    def __repr__(self):
        return str(self.e)

class triangle:
    def __init__(self, a, b, c):
        #упорядочить?
        #сортировка по x:
        if a.x>b.x:
            a, b = b, a
        if a.x>c.x:
            a, c = c, a
        if b.x>c.x:
            b, c = c, b
        
        if a.x==b.x:
            if a.y>b.y:
                a, b = b, a

        #обход по часовой стрелке
        self.eb = edge(a, c)
        self.ec = edge(a, b)
        if self.eb.v.vprod(self.ec.v)<0:
            b, c = c, b
    
        self.a = a
        self.b = b
        self.c = c
        self.ea = edge(b, c)
        self.eb = edge(c, a)
        self.ec = edge(a, b)

    def S(self,):
        return abs( (self.b.x-self.a.x)*(self.c.y-self.a.y) - (self.c.x-self.a.x)*(self.b.y-self.a.y) )/2

    def V(self,):
        def det(a, b, c):
            r1 = a.x*b.y*c.z
            r2 = a.y*b.z*c.x
            r3 = a.z*b.x*c.y
            p1 = a.z*b.y*c.x
            p2 = a.x*b.z*c.y
            p3 = a.y*b.x*c.z
            return (r1+r2+r3-p1-p2-p3)
        s = self.S()
        h = [self.a, self.b, self.c,]
        h.sort(key = lambda q: q.z)
        h10 = dot(h[1].x, h[1].y, h[0].z)
        h20 = dot(h[2].x, h[2].y, h[0].z)

        h1h0 = vect(h[1], h[0])
        h1h2 = vect(h[1], h[2])
        h1h20= vect(h[1], h20)
        h1h10= vect(h[1], h10)
        vh = abs( det(h1h0, h1h2, h1h20) )
        vl = abs( det(h1h0, h1h10, h1h20) )
        
        v =s*h[0].z + (vh+vl)/6
        #врооде корректный объём. считаю 2 смешаных произведения,
        #рассекая хреновину по средней по высоте точке на 2 тетраэдера
        return v

    def dot_inside(self, d):
        if d == self.a:
            return False
        elif d == self.b:
            return False
        elif d == self.c:
            return False
        vad = vect(self.a, d)
        vbd = vect(self.b, d)
        vcd = vect(self.c, d)
        return ( (self.ec.v.vprod(vad)<=0) and (self.ea.v.vprod(vbd)<=0) and (self.eb.v.vprod(vcd)<=0) )

    def __eq__(self, d):
        return (self.a==d.a)and(self.b==d.b)and(self.c==d.c)

    def __hash__(self):
        return hash((self.a, self.b, self.c))

    def __repr__(self):
        return 'A'+str((self.a.x, self.a.y))+' B'+str((self.b.x, self.b.y))+' C'+str((self.c.x, self.c.y))

class scan:
    #попробовать перегрузить классы и добавить в них поле "номер вершины". или как то так
    # если будет медленно - попробовать https://docs.python.org/3/library/heapq.html
    def __init__(self, x, y, z):
        self.x = x.copy()
        self.y = y.copy()
        self.z = z.copy()
        self.used = [False]*len(self.x)
        self.points = [dot(d[0],d[1],d[2]) for d in zip(self.x, self.y, self.z)]
        self.points.sort(key = lambda d: d.x)
        self.edges = set()
        self.edges_sorted = []
        self.triangles = set()
        self.tmp_results = {}
        
    def add_edge(self, e):
        if e in self.edges:
            return
        self.edges.add(e)
        self.edges_sorted.append(e)
        self.edges_sorted.sort(key = lambda e: e.e[0].x)

    def start_edge(self, b = 0, e = 1):
        self.edges = set()
        self.edges_sorted = []
        self.add_edge( edge(self.points[b], self.points[e]) )

    #проверка на пересечение отрезков внутренними частями, без учёта точек на краю
    def check_inner_crossing(self, j, left_bound, right_bound):
        #j == ab, e == cd
        first_seg = j.v
        for e in self.edges_sorted[left_bound : right_bound]:
            second_seg = vect(e.e[0], e.e[1])
            
            ac = vect(j.e[0], e.e[0])
            ad = vect(j.e[0], e.e[1])

            ca = vect(e.e[0], j.e[0])
            cb = vect(e.e[0], j.e[1])
            #скалярное *
            #векторное ^
            S_with_j = first_seg.vprod(ac)*first_seg.vprod(ad)
            S_with_e = second_seg.vprod(ca)*second_seg.vprod(cb)
            if (S_with_j < 0) and  (S_with_e < 0):
                return True
        return False

    def check_dots_inside(self, t, left_bound, right_bound):
        for point in self.points[left_bound : right_bound]:
            if t.dot_inside(point):
                return True
        return False

    def find_edges_bounds(self, l_c, r_c):
        #индексы границ в массиве
        lb = 0
        rb = len(self.edges_sorted)

        while rb-lb>1:
            m = (rb+lb)//2
            if self.edges_sorted[m].e[0].x<l_c:
                lb=m
            else:
                rb=m
        l_b = lb

        rb = len(self.edges_sorted)
        while rb-lb>1:
            m = (rb+lb)//2
            if self.edges_sorted[m].e[0].x<=r_c:
                lb=m
            else:
                rb=m
        r_b = rb+1
        return l_b, r_b

    def find_points_bounds(self, l_c, r_c):
        #индексы границ в массиве
        lb = 0
        rb = len(self.points)

        while rb-lb>1:
            m = (rb+lb)//2
            if self.points[m].x<l_c:
                lb=m
            else:
                rb=m
        l_b = lb

        rb = len(self.points)
        while rb-lb>1:
            m = (rb+lb)//2
            if self.points[m].x<=r_c:
                lb=m
            else:
                rb=m
        r_b = rb+1
        return l_b, r_b

    #ищу ближайшую, но не образующую уже имеющийся треугольник
    def get_closest_free_point_to_edge(self, j, min_len):
        bound = pow(min_len, 0.5)
        if j not in self.tmp_results:
            self.tmp_results[j]=[]
            #граничные координаты
            l_c = j.e[0].x - bound
            r_c = j.e[1].x + bound
            left_bound, right_bound = self.find_points_bounds(l_c, r_c)
            for d in self.points[left_bound : right_bound]:
                if d in j.e:
                    continue
                ea = edge(j.e[0], d)
                eb = edge(j.e[1], d)
                t = triangle(j.e[0], j.e[1], d)
                if t in self.triangles:
                    continue
                #надо ещё проверить что рёбра не на одной прямой!
                if (ea.v.vprod(eb.v))==0:
                    continue
                #тут проверка на пересечение
                if self.check_inner_crossing(ea, *self.find_edges_bounds(ea.e[0].x-bound, ea.e[1].x)):
                    continue
                if self.check_inner_crossing(eb, *self.find_edges_bounds(eb.e[0].x-bound, eb.e[1].x)):
                    continue
                #тут проверка на включение точки внутрь треугольника
                #с учётом рёбер, без учёта вершин
                if self.check_dots_inside(t, *self.find_points_bounds(t.a.x, t.c.x)):
                    continue
                distance = j.L2norm(d)
                #кажется такой треугольник может подойти. Добавляем
                if distance >= min_len:
                    continue
                self.tmp_results[j].append((d, distance, t.S()))
                #min_len = distance
            #Наконец - сортировка. если есть альетрнатива с большей площадью - выбрать её
            self.tmp_results[j].sort(key = lambda d: d[2],reverse=False)
            self.tmp_results[j].sort(key = lambda d: d[1],reverse=True)
            if self.tmp_results[j]:
                result_dot, result_len, _ = self.tmp_results[j][-1]
            else:
                result_dot = None
                result_len = min_len
            return result_dot, result_len

        result_dot = None
        result_len = min_len
        while self.tmp_results[j]:
            result_dot, result_len, max_s = self.tmp_results[j][-1]
            ea = edge(j.e[0], result_dot)
            eb = edge(j.e[1], result_dot)
            if self.check_inner_crossing(ea, *self.find_edges_bounds(ea.e[0].x-bound, ea.e[1].x)) or self.check_inner_crossing(eb, *self.find_edges_bounds(eb.e[0].x-bound, eb.e[1].x)):
                result_dot = None
                result_len = min_len
                self.tmp_results[j].pop()
                continue
            break
        return result_dot, min_len

    #min_len - сумма квадратов длинн от краев ребра до точки!!
    def get_next_triangle(self, min_len = 2*2):
        good_edge = None
        good_dot = None
        tmp_min_len = min_len
        #от этого цикла можно избавиться, если закинуть все возможные треугольники в 1 лист и отсортировать
        for j in self.edges:
            d, tmp_min_len = self.get_closest_free_point_to_edge(j, min_len)
            if not(d is None) and tmp_min_len<=min_len:
                good_dot = d
                good_edge = j
                good_len = tmp_min_len
        if good_dot:
            print('d=', good_dot, 'j=', good_edge,)# "list:", self.tmp_results[good_edge])
            self.tmp_results[good_edge].pop()
        return good_edge, good_dot

    def add_triangle(self, j, d):
        self.add_edge( edge(j.e[0], d) )
        self.add_edge( edge(j.e[1], d) )
        self.triangles.add( triangle(j.e[0], j.e[1], d) )

    def triangulate(self,):
        self.start_edge()
        for i in range(len(self.x)**2):
            j, d = self.get_next_triangle(8)
            if not(d):
                break
            self.add_triangle(j, d)

    def calc_volume(self,):
        total_fv = 0
        for t in self.triangles:
            v = t.V()
            #print(t, "volume=", v)
            total_fv += v
        self.V = total_fv

if __name__ == '__main__':
    q = [1, 3, 5, 11, 11, 11, 12, 12, 13, 13, 15, 17, 17, 21, 21,]
    print('len:', len(q))
    num = 11

    lb = 0
    rb = len(q)
    while rb-lb>1:
        print(lb, rb)
        m = (rb+lb)//2
        if q[m]<=num:
            lb=m
        else:
            rb=m
    
    print(lb, rb,)
    print(q[lb:rb+1])
