CREATE TABLE IF NOT EXISTS messages (
    id UUID NOT NULL DEFAULT uuidv7(),
    telegram_user_id BIGINT NOT NULL,
    telegram_user_name VARCHAR(32) DEFAULT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

COMMENT ON TABLE messages IS 'Таблица сообщений отправленных в telegram бот';
COMMENT ON COLUMN messages.id IS 'Идентификатор сообщения';
COMMENT ON COLUMN messages.telegram_user_id IS 'Идентификатор telegram пользователя';
COMMENT ON COLUMN messages.telegram_user_name IS 'Имя telegram пользователя';
COMMENT ON COLUMN messages.text IS 'Текст сообщения';
COMMENT ON COLUMN messages.created_at IS 'дата и время создания записи';

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'sunt', 'quia et suscipit
suscipit recusandae consequuntur expedita et cum
reprehenderit molestiae ut ut quas totam
nostrum rerum est autem sunt rem eveniet architecto');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'molestias', 'et iusto sed quo iure
voluptatem occaecati omnis eligendi aut ad
voluptatem doloribus vel accusantium quis pariatur
molestiae porro eius odio et labore et velit aut');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'eum', 'ullam et saepe reiciendis voluptatem adipisci
sit amet autem assumenda provident rerum culpa
quis hic commodi nesciunt rem tenetur doloremque ipsam iure
quis sunt voluptatem rerum illo velit');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'nesciunt', 'repudiandae veniam quaerat sunt sed
alias aut fugiat sit autem sed est
voluptatem omnis possimus esse voluptatibus quis
est aut tenetur dolor neque');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'dolorem', 'ut aspernatur corporis harum nihil quis provident sequi
mollitia nobis aliquid molestiae
perspiciatis et ea nemo ab reprehenderit accusantium quas
voluptate dolores velit et doloremque molestiae');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'magnam', 'dolore placeat quibusdam ea quo vitae
magni quis enim qui quis quo nemo aut saepe
quidem repellat excepturi ut quia
sunt ut sequi eos ea sed quas');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'dolorem', 'dignissimos aperiam dolorem qui eum
facilis quibusdam animi sint suscipit qui sint possimus cum
quaerat magni maiores excepturi
ipsam ut commodi dolor voluptatum modi aut vitae');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'nesciunt', 'consectetur animi nesciunt iure dolore
enim quia ad
veniam autem ut quam aut nobis
et est aut quod aut provident voluptas autem voluptas');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'optio', 'quo et expedita modi cum officia vel magni
doloribus qui repudiandae
vero nisi sit
quos veniam quod sed accusamus veritatis error');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'etea', 'delectus reiciendis molestiae occaecati non minima eveniet qui voluptatibus
accusamus in eum beatae sit
vel qui neque voluptates ut commodi qui incidunt
ut animi commodi');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'quibusdam', 'itaque id aut magnam
praesentium quia et ea odit et ea voluptas et
sapiente quia nihil amet occaecati quia id voluptatem
incidunt ea est distinctio odio');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'dolorum', 'aut dicta possimus sint mollitia voluptas commodi quo doloremque
iste corrupti reiciendis voluptatem eius rerum
sit cumque quod eligendi laborum minima
perferendis recusandae assumenda consectetur porro architecto ipsum ipsam');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'voluptatem', 'fuga et accusamus dolorum perferendis illo voluptas
non doloremque neque facere
ad qui dolorum molestiae beatae
sed aut voluptas totam sit illum');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'eveniet', 'reprehenderit quos placeat
velit minima officia dolores impedit repudiandae molestiae nam
voluptas recusandae quis delectus
officiis harum fugiat vitae');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'sint', 'suscipit nam nisi quo aperiam aut
asperiores eos fugit maiores voluptatibus quia
voluptatem quis ullam qui in alias quia est
consequatur magni mollitia accusamus ea nisi voluptate dicta');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'fugit', 'eos voluptas et aut odit natus earum
aspernatur fuga molestiae ullam
deserunt ratione qui eos
qui nihil ratione nemo velit ut aut id quo');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'voluptate', 'eveniet quo quis
laborum totam consequatur non dolor
ut et est repudiandae
est voluptatem vel debitis et magnam');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'adipisci', 'illum quis cupiditate provident sit magnam
ea sed aut omnis
veniam maiores ullam consequatur atque
adipisci quo iste expedita sit quos voluptas');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'doloribus', 'qui consequuntur ducimus possimus quisquam amet similique
suscipit porro ipsam amet
eos veritatis officiis exercitationem vel fugit aut necessitatibus totam
omnis rerum consequatur expedita quidem cumque explicabo');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'asperiores', 'repellat aliquid praesentium dolorem quo
sed totam minus non itaque
nihil labore molestiae sunt dolor eveniet hic recusandae veniam
tempora et tenetur expedita sunt');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'dolor', 'eos qui et ipsum ipsam suscipit aut
sed omnis non odio
expedita earum mollitia molestiae aut atque rem suscipit
nam impedit esse');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'maxime', 'veritatis unde neque eligendi
quae quod architecto quo neque vitae
est illo sit tempora doloremque fugit quod
et et vel beatae sequi ullam sed tenetur perspiciatis');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'autem', 'enim et ex nulla
omnis voluptas quia qui
voluptatem consequatur numquam aliquam sunt
totam recusandae id dignissimos aut sed asperiores deserunt');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'rem', 'ullam consequatur ut
omnis quis sit vel consequuntur
ipsa eligendi ipsum molestiae et omnis error nostrum
molestiae illo tempore quia et distinctio');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'est', 'similique esse doloribus nihil accusamus
omnis dolorem fuga consequuntur reprehenderit fugit recusandae temporibus
perspiciatis cum ut laudantium
omnis aut molestiae vel vero');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'quasi', 'eum sed dolores ipsam sint possimus debitis occaecati
debitis qui qui et
ut placeat enim earum aut odit facilis
consequatur suscipit necessitatibus rerum sed inventore temporibus consequatur');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'delectus', 'non et quaerat ex quae ad maiores
maiores recusandae totam aut blanditiis mollitia quas illo
ut voluptatibus voluptatem
similique nostrum eum');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'iusto', 'odit magnam ut saepe sed non qui
tempora atque nihil
accusamus illum doloribus illo dolor
eligendi repudiandae odit magni similique sed cum maiores');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'aquo', 'alias dolor cumque
impedit blanditiis non eveniet odio maxime
blanditiis amet eius quis tempora quia autem rem
a provident perspiciatis quia');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'ullam', 'debitis eius sed quibusdam non quis consectetur vitae
impedit ut qui consequatur sed aut in
quidem sit nostrum et maiores adipisci atque
quaerat voluptatem adipisci repudiandae');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'doloremque', 'deserunt eos nobis asperiores et hic
est debitis repellat molestiae optio
nihil ratione ut eos beatae quibusdam distinctio maiores
earum voluptates et aut adipisci ea maiores voluptas maxime');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'qui', 'rerum ut et numquam laborum odit est sit
id qui sint in
quasi tenetur tempore aperiam et quaerat qui in
rerum officiis sequi cumque quod');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'magnam', 'ea velit perferendis earum ut voluptatem voluptate itaque iusto
totam pariatur in
nemo voluptatem voluptatem autem magni tempora minima in
est distinctio qui assumenda accusamus dignissimos officia nesciunt nobis');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'id', 'nisi error delectus possimus ut eligendi vitae
placeat eos harum cupiditate facilis reprehenderit voluptatem beatae
modi ducimus quo illum voluptas eligendi
et nobis quia fugit');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'fuga', 'ad mollitia et omnis minus architecto odit
voluptas doloremque maxime aut non ipsa qui alias veniam
blanditiis culpa aut quia nihil cumque facere et occaecati
qui aspernatur quia eaque ut aperiam inventore');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'provident', 'debitis et eaque non officia sed nesciunt pariatur vel
voluptatem iste vero et ea
numquam aut expedita ipsum nulla in
voluptates omnis consequatur aut enim officiis in quam qui');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'explicabo', 'animi esse sit aut sit nesciunt assumenda eum voluptas
quia voluptatibus provident quia necessitatibus ea
rerum repudiandae quia voluptatem delectus fugit aut id quia
ratione optio eos iusto veniam iure');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'eos', 'corporis rerum ducimus vel eum accusantium
maxime aspernatur a porro possimus iste omnis
est in deleniti asperiores fuga aut
voluptas sapiente vel dolore minus voluptatem incidunt ex');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'enim', 'ut voluptatum aliquid illo tenetur nemo sequi quo facilis
ipsum rem optio mollitia quas
voluptatem eum voluptas qui
unde omnis voluptatem iure quasi maxime voluptas nam');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'non', 'molestias id nostrum
excepturi molestiae dolore omnis repellendus quaerat saepe
consectetur iste quaerat tenetur asperiores accusamus ex ut
nam quidem est ducimus sunt debitis saepe');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'commodi', 'odio fugit voluptatum ducimus earum autem est incidunt voluptatem
odit reiciendis aliquam sunt sequi nulla dolorem
non facere repellendus voluptates quia
ratione harum vitae ut');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'eligendi', 'similique fugit est
illum et dolorum harum et voluptate eaque quidem
exercitationem quos nam commodi possimus cum odio nihil nulla
dolorum exercitationem magnam ex et a et distinctio debitis');
